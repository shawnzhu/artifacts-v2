package main

import (
	"os"

	"net/http"

	"github.com/travis-ci/artifacts-v2/router"
	"github.com/urfave/cli"
)

func server(c *cli.Context) error {
	return http.ListenAndServe(c.String("server-addr"), router.Routes())
}

func app() *cli.App {
	app := cli.NewApp()

	app.Name = "travis-artifacts"
	app.Usage = "start a server to access artifacts"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server-addr",
			Usage:  "server address",
			Value:  ":8080",
			EnvVar: "SERVER_ADDR",
		},
	}

	app.Action = server

	return app
}

func main() {
	app().Run(os.Args)
}
