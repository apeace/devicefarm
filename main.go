package main

import (
	"flag"
	"github.com/ride/devicefarm/build"
	"log"
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

	err = build.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
