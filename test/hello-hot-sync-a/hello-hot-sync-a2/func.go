package main

import (
	"context"
	fdk "github.com/fnproject/fdk-go"
	"io"
	"time"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {

	// Always forces a timeout

	out.Write([]byte("FORCETIMEOUT")) // Used by tests to check the function completed OK (async only
	time.Sleep(time.Duration(60) * time.Second)

	// don't expect to get this far
	out.Write([]byte("Hello from hello-hot-sync-a/hot-hot-sync-a2"))
	out.Write([]byte("COMPLETEDOK")) // Used by tests to check the function completed OK

}
