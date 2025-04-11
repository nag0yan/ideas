package main

import (
	"fmt"
	"time"

	"github.com/nag0yan/ideas/logic"
)

type Timer struct{}

func (t Timer) Sleep(d time.Duration) {
	time.Sleep(d)
}

func main() {
	fmt.Println("Starting...")
	logic.Retry(CauseError, 3, time.Millisecond*1000, &Timer{})
	fmt.Println("Finished")
}

func CauseError() (int, error) {
	return -1, fmt.Errorf("an error occurred")
}
