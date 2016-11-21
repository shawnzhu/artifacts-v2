package main

import (
	"net/http"

	"github.com/travis-ci/artifacts-v2/router"
)

func main() {
	http.ListenAndServe(":8080", router.Routes())
}
