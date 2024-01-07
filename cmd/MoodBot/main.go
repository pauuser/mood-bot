package main

import (
	"flag"
	"fmt"
	"log"
	"pauuser/mood-bot/internal/app"
)

func main() {
	application := app.App{}
	pathToConfig := flag.String("path-to-config", "./config", "path to dir with config")
	configFileName := flag.String("config-file-name", "config", "config file name")
	flag.Parse()
	fmt.Println("path to config =", *pathToConfig, "config file name =", *configFileName)

	err := application.ParseConfig(*pathToConfig, *configFileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	application.Run()

}
