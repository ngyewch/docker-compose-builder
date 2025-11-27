package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
)

var (
	version string

	app = &cli.Command{
		Name:           "docker-compose-builder",
		Usage:          "docker compose builder",
		DefaultCommand: "build",
		Commands: []*cli.Command{
			{
				Name:   "build",
				Usage:  "build",
				Action: doBuild,
			},
		},
	}
)

func main() {
	if version == "" {
		buildInfo, _ := debug.ReadBuildInfo()
		if buildInfo != nil {
			version = buildInfo.Main.Version
		}
	}
	app.Version = version

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		oopsErr, ok := oops.AsOops(err)
		if ok {
			fmt.Println(oopsErr.Stacktrace())
			return
		}
		log.Fatal(err)
	}
}
