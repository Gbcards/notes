package main

import (
	"notes-service/internal/platform"

	"go.uber.org/fx"
)

func main() {
	fx.New(platform.Modules()...).Run()
}
