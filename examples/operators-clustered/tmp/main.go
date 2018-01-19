package main

import (
	"context"

	"github.com/fnproject/fn/api/server"
		_ "github.com/fnproject/ext-statsapi/stats"
)

func main() {
	ctx := context.Background()
	funcServer := server.NewFromEnv(ctx)
		funcServer.AddExtensionByName("github.com/fnproject/ext-statsapi/stats")
	funcServer.Start(ctx)
}
