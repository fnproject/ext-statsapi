package main

import (
	"time"
)

func main() {
	var sleepTimeMilliseconds int = 200
	time.Sleep(time.Duration(sleepTimeMilliseconds) * time.Millisecond)
}