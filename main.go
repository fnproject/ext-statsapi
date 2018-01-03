package main

// This file is used only by developers who wish to build a local fn server executable
// See examples/developers/README.md for more information

import (
	"context"
	"fmt"

	"github.com/fnproject/ext-statsapi/stats"
	"github.com/fnproject/fn/api/server"
)

func main() {
	ctx := context.Background()
	funcServer := server.NewFromEnv(ctx)
	funcServer.AddExtensionByName(stats.StatisticsExtensionName())

	fmt.Println("Starting fn server with metrics API extensions")
	funcServer.Start(ctx)
}
