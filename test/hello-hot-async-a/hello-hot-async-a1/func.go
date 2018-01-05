package main

import (
	"context"
	fdk "github.com/fnproject/fdk-go"
	"io"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {

	out.Write([]byte("Hello from hello-hot-async-a/hot-hot-async-a1"))
	out.Write([]byte("COMPLETEDOK")) // Used by tests to check the function completed OK

}
