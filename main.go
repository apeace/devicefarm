package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/codegangsta/cli"
	"github.com/ride/devicefarm/awsutil"
	"github.com/ride/devicefarm/build"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

var currentUser *user.User
var defaultAwsConfigFile string

func init() {
	var err error
	currentUser, err = user.Current()
	if err != nil {
		log.Fatalln("Could not get current user")
	}
	defaultAwsConfigFile = filepath.Join(currentUser.HomeDir, ".devicefarm.json")
}

func main() {
	app := cli.NewApp()
	app.Name = "devicefarm"
	app.Usage = "Run UI tests in AWS Device Farm"

	// these flags are used for anything which needs the context of
	// a build directory
	buildFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Usage: "Working directory of the build",
			Value: ".",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "Config file relative to working directory (or absolute path)",
			Value: "devicefarm.yml",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "run",
			Usage:     "Create test run based on YAML config",
			ArgsUsage: " ",
			Action:    commandRun,
			Flags:     buildFlags,
		},
		{
			Name:      "build",
			Usage:     "Run local build based on YAML config",
			ArgsUsage: " ",
			Action:    commandBuild,
			Flags:     buildFlags,
		},
		{
			Name:      "devicepools",
			Usage:     "Sync devicepools with your YAML config",
			ArgsUsage: " ",
			Action:    commandDevicePools,
			Flags:     buildFlags,
		},
		{
			Name:      "devices",
			Usage:     "Search device farm devices",
			ArgsUsage: "[search]",
			Action:    commandDevices,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "android",
					Usage: "Filter to only Android devices",
				},
				cli.BoolFlag{
					Name:  "ios",
					Usage: "Filter to only iOS devices",
				},
			},
		},
	}

	app.Run(os.Args)
}

func commandRun(c *cli.Context) {
	commandBuild(c)
	pool := getDevicePool(c)
	build := getBuild(c)
	client := getClient()
	//lookup := client.
	apk := filepath.Join(build.Dir, build.Manifest.Android.Apk)
	apkInstrumentation := filepath.Join(build.Dir, build.Manifest.Android.ApkInstrumentation)
	runArn, err := client.CreateRun(build.Config.Arn, *pool.Arn, apk, apkInstrumentation)
	if err != nil {
		log.Fatalln(err)
	}
	re := regexp.MustCompile("run:([^/]+)/([^/]+)")
	parts := re.FindStringSubmatch(runArn)
	log.Printf("https://us-west-2.console.aws.amazon.com/devicefarm/home?region=us-west-2#/projects/%s/runs/%s\n", parts[1], parts[2])
}

func commandBuild(c *cli.Context) {
	build := getBuild(c)
	log.Println(">> Running build... (silencing output)")
	err := build.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(">> Build complete")
}

func getDevicePool(c *cli.Context) *devicefarm.DevicePool {
	build := getBuild(c)
	client := getClient()

	poolDefs, err := build.Config.FlatDevicePoolDefinitions()
	if err != nil {
		log.Fatalln(err)
	}

	pools, err := client.ListDevicePools(build.Config.Arn)
	if err != nil {
		log.Fatalln(err)
	}

	poolName := build.Manifest.DevicePool
	def, ok := poolDefs[poolName]
	if !ok {
		log.Fatalln("Device Pool not defined: " + poolName)
	}
	log.Printf(">> Device Pool: %s (%d devices)\n", poolName, len(def))

	remoteName := "df:" + build.Branch + ":" + poolName
	var matchingPool *devicefarm.DevicePool
	for _, pool := range pools {
		if *pool.Name == remoteName {
			matchingPool = pool
		}
	}

	if matchingPool == nil {
		log.Println("...creating")
		matchingPool, err = client.CreateDevicePool(build.Config.Arn, remoteName, def)
		if err != nil {
			log.Fatalln(err)
		}
	}

	matches, err := client.DevicePoolMatches(matchingPool, def)
	if err != nil {
		log.Fatalln(err)
	}
	if !matches {
		log.Println("...updating")
		matchingPool, err = client.UpdateDevicePool(matchingPool, def)
	}

	return matchingPool
}

func commandDevicePools(c *cli.Context) {
	getDevicePool(c)
}

func commandDevices(c *cli.Context) {
	client := getClient()
	search := ""
	if c.NArg() > 0 {
		search = c.Args()[0]
	}
	androidOnly := c.Bool("android")
	iosOnly := c.Bool("ios")
	if androidOnly && iosOnly {
		log.Fatalln("Cannot use both --android and --ios")
	}
	devices, err := client.SearchDevices(search, androidOnly, iosOnly)
	if err != nil {
		log.Fatalln(err)
	}
	for _, device := range devices {
		fmt.Println(*device.Name)
	}
}

func findCreds() *credentials.Credentials {
	ok, creds := awsutil.CredsFromEnv()
	if !ok {
		ok, creds = awsutil.CredsFromFile(defaultAwsConfigFile)
	}
	if !ok {
		log.Fatalln("Could not find AWS credentials")
	}
	return creds
}

func getClient() *awsutil.DeviceFarm {
	creds := findCreds()
	return awsutil.NewClient(creds)
}

var cachedBuild *build.Build

func getBuild(c *cli.Context) *build.Build {
	if cachedBuild != nil {
		return cachedBuild
	}

	dir := c.String("dir")
	configFile := c.String("config")

	log.Println(">> Dir: " + dir)
	log.Println(">> Config: " + configFile)

	if len(dir) > 1 && dir[:2] == "~/" {
		dir = filepath.Join(currentUser.HomeDir, dir[2:])
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln(err)
	}
	absConfigFile := configFile
	if !filepath.IsAbs(configFile) {
		absConfigFile = filepath.Join(absDir, configFile)
	}

	build, err := build.New(absDir, absConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(">> Branch: " + build.Branch)

	cachedBuild = build

	return build
}
