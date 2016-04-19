package main

import (
	"flag"
	"fmt"
	"github.com/ride/devicefarm/config"
)

func main() {
	flag.Parse()
	configFile := flag.Arg(0)
	_, err := config.New(configFile)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Your config is valid")
	}
}
