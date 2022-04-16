package main

import (
	"net/http"

	stats "github.com/hi-bridge-9/mlb-players-stats"
)

func main() {
	var w http.ResponseWriter
	var r *http.Request

	stats.FunctionTest(w, r)

}
