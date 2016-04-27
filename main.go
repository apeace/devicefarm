package main

import (
	"flag"
	"github.com/ride/devicefarm/build"
	"log"
	"os/user"
	"path"
	"path/filepath"
)

func main() {
	var dir string
	var configFile string
	flag.StringVar(&dir, "dir", ".",
		"Working directory of the build")
	flag.StringVar(&configFile, "file", "devicefarm.yml",
		"Config file relative to `dir`, or an absolute path")
	flag.Parse()

	log.Println(">> Dir: " + dir)
	log.Println(">> Config: " + configFile)

	if dir[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalln("Could not get current user")
		}
		dir = path.Join(usr.HomeDir, dir[2:])
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

	log.Println(">> Running build... (silencing output)")
	err = build.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(">> Build complete")
}
