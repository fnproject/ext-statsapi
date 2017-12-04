package main

import (
	"context"
	"fmt"
	"github.com/fnproject/ext-metrics/statistics"
	"github.com/fnproject/fn/api/server"
)

func main() {
	ctx := context.Background()
	funcServer := server.NewFromEnv(ctx)

	// Add the metrics API extension
	statistics.AddEndpoints(funcServer)

	fmt.Println("Starting fn server with metrics API extensions")
	funcServer.Start(ctx)
}
