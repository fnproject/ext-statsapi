package main

import (
	"context"
	"fmt"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/ext-metrics/handlers"
)

func main() {
	ctx := context.Background()
	funcServer := server.NewFromEnv(ctx)

	// Add the metrics API extension
	handlers.AddEndpoints(funcServer)

	fmt.Println("Starting fn server with metrics API extensions")
	funcServer.Start(ctx)
}
