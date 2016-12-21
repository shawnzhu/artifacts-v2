package main

import (
	"bytes"
	"fmt"
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

		g.Before(func() {
			app.Action = func(c *cli.Context) error {
				fmt.Printf("Hello World")
				return nil
			}

			app.Writer = w

			err := app.Run([]string{"travis-artifacts", "-h"})
			g.Assert(err).Equal(nil)
		})

		g.It("shows usage", func() {
			g.Assert(bytes.Contains(w.written, []byte(app.Usage)))
		})

		g.It("supports parameter server-addr", func() {
			g.Assert(bytes.Contains(w.written, []byte("server-addr")))
		})
	})
}
