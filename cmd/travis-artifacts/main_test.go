package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/urfave/cli"

	. "github.com/franela/goblin"
)

type mockWriter struct {
	written []byte
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	if w.written == nil {
		w.written = p
	} else {
		w.written = append(w.written, p...)
	}

	return len(p), nil
}

func (w *mockWriter) GetWritten() (b []byte) {
	return w.written
}

func TestCli(t *testing.T) {
	g := Goblin(t)

	g.Describe("Run ./travis-artifacts -h", func() {
		w := &mockWriter{}
		app := app()
		var output string

		g.Before(func() {
			app.Action = func(c *cli.Context) error {
				fmt.Printf("Hello World")
				return nil
			}

			app.Writer = w

			err := app.Run([]string{"travis-artifacts", "-h"})
			output = string(w.GetWritten())
			g.Assert(err).Equal(nil)
		})

		g.It("shows usage", func() {
			g.Assert(strings.Contains(output, app.Usage)).Equal(true)
		})

		params := [...]string{
			"server-addr",
			"server-cert",
			"server-key",
			"db-url",
			"jwt-public-key",
		}

		for _, param := range params {
			g.It(fmt.Sprintf("supports parameter %s", param), func() {
				g.Assert(strings.Contains(output, param)).Equal(true)
			})
		}
	})
}
