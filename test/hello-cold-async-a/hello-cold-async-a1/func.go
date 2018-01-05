package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {

	// Get query parameters forcepanic and forcetimeout
	// if both forcepanic and forcetimeout are set then function will panic

	s := os.Getenv("FN_REQUEST_URL")
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
			fmt.Println("FORCETIMEOUT") // Used by tests to check the function completed OK (async only
			time.Sleep(time.Duration(60) * time.Second)
		}
	}

	fmt.Println("Hello from hello-cold-async-a/hello-cold-async-a1")
	fmt.Println("COMPLETEDOK") // Used by tests to check the function completed OK
}
