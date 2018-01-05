package main

import (
	"context"
	"errors"
	fdk "github.com/fnproject/fdk-go"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {

	// Looks as if "FN_REQUEST_URL" is not set in hot functions, so forcepanic and forcetimeout are ignored

	// Get query parameters forcepanic and forcetimeout
	// if both forcepanic and forcetimeout are set then function will panic

	s := os.Getenv("FN_REQUEST_URL")
	out.Write([]byte("<" + s + ">"))

	a := os.Getenv("FN_APP_NAME")
	out.Write([]byte("<" + a + ">"))

	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)

	// The query parameter "forcepanic" can be specified to true or false
	// If true this function will panic
	forcePanicParam := m.Get("forcepanic")
	if forcePanicParam != "" {
		forcePanic, err := strconv.ParseBool(forcePanicParam)
		if forcePanic || err != nil {
			panic(errors.New("FORCEPANIC"))
		}
	}

	// The query parameter "forceTimeout"  can be specified to true or false
	// if true this function will timeout
	forceTimeoutParam := m.Get("forcetimeout")
	if forceTimeoutParam != "" {
		forceTimeout, err := strconv.ParseBool(forceTimeoutParam)
		if forceTimeout || err != nil {
			out.Write([]byte("FORCETIMEOUT")) // Used by tests to check the function completed OK (async only
			time.Sleep(time.Duration(60) * time.Second)
		}
	}

	out.Write([]byte("Hello from hello-hot-sync-a/hot-hot-sync-a1"))
	out.Write([]byte("COMPLETEDOK")) // Used by tests to check the function completed OK

}
