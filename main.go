package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/codegangsta/cli"
	"github.com/ride/devicefarm/awsutil"
	"github.com/ride/devicefarm/build"
	"log"
	"os"
	"os/user"
	"path/filepath"
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
			Name:      "build",
			Usage:     "Run build based on YAML config",
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

func commandBuild(c *cli.Context) {
	build := getBuild(c)
	log.Println(">> Running build... (silencing output)")
	err := build.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(">> Build complete")
}

func commandDevicePools(c *cli.Context) {
	build := getBuild(c)
	log.Println(build.Dir)
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
	devices, err := client.ListDevices(search, androidOnly, iosOnly)
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

func getBuild(c *cli.Context) *build.Build {
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

	return build
}
