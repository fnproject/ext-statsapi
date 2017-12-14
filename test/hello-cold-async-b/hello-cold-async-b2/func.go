package main

import (
	"fmt"
	"time"
)

func main() {

	var sleepTimeMilliseconds int = 200
	time.Sleep(time.Duration(sleepTimeMilliseconds) * time.Millisecond)
	fmt.Println("Hello from hello-cold-async-b/hello-cold-async-b2")
}
