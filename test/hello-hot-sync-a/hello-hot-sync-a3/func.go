package main

import (
	"context"
	"errors"
	fdk "github.com/fnproject/fdk-go"
	"io"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {

	// Always panics
	panic(errors.New("FORCEPANIC"))

}
