package main

import (
	"context"
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/mariospas/gqlgenc/clientgen"
	"github.com/mariospas/gqlgenc/clientgenv2"
	"github.com/mariospas/gqlgenc/config"
	"github.com/mariospas/gqlgenc/generator"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err.Error())
		os.Exit(2)
	}

	clientGen := api.AddPlugin(clientgen.New(cfg.Query, cfg.Client, cfg.Generate))
	if cfg.Generate != nil {
		if cfg.Generate.ClientV2 {
			clientGen = api.AddPlugin(clientgenv2.New(cfg.Query, cfg.Client, cfg.Generate))
		}
	}

	if err := generator.Generate(ctx, cfg, clientGen); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err.Error())
		os.Exit(4)
	}
}
