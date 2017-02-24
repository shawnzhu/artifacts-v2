package main

import (
	"os"

	"net/http"

	"github.com/travis-ci/artifacts-v2/router"
	"github.com/urfave/cli"
)

func server(c *cli.Context) error {
	handler := router.Routes(c)

	if c.String("server-cert") == "" {
		return http.ListenAndServe(c.String("server-addr"), handler)
	}

	// TLS support
	return http.ListenAndServeTLS(
		c.String("server-addr"),
		c.String("server-cert"),
		c.String("server-key"),
		handler,
	)
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
		cli.StringFlag{
			Name:   "server-cert",
			Usage:  "server TLS cert",
			EnvVar: "SERVER_CERT",
		},
		cli.StringFlag{
			Name:   "server-key",
			Usage:  "server TLS key",
			EnvVar: "SERVER_KEY",
		},
		cli.StringFlag{
			Name:   "db-url",
			Usage:  "database URL",
			EnvVar: "DB_URL",
		},
		cli.StringFlag{
			Name:   "jwt-public-key",
			Usage:  "RSA public key for JWT",
			EnvVar: "JWT_PUBLIC_KEY",
		},
	}

	app.Action = server

	return app
}

func main() {
	app().Run(os.Args)
}
