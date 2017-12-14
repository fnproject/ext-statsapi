package main

import (
	"context"
	"encoding/json"
	"fmt"
	fdk "github.com/fnproject/fdk-go"
	"io"
	"time"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {

	doStuff()

	var person struct {
		Name string `json:"name"`
	}
	json.NewDecoder(in).Decode(&person)
	if person.Name == "" {
		person.Name = "world"
	}

	msg := struct {
		Msg string `json:"msg"`
	}{
		Msg: fmt.Sprintf("Hello %s!\n", person.Name),
	}

	json.NewEncoder(out).Encode(&msg)
}

func doStuff() {
	var sleepTimeMilliseconds int = 200
	time.Sleep(time.Duration(sleepTimeMilliseconds) * time.Millisecond)
}
