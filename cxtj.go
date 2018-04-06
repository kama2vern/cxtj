package main

import (
	"os"
	"runtime"

	"./logger"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "cxtj"
	app.Version = "0.0.1" // TODO: dynamic version & commit hash
	app.Usage = "A CLI tool for conversion from xlsx to json"
	app.Author = "kama2vern"
	app.Commands = Commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: "",
			Usage: "Config file path",
		},
	}

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	err := app.Run(os.Args)
	if err != nil {
		logger.Log("error", err.Error())
		os.Exit(1)
	}
}
