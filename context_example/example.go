package context_example

import (
	"context"
	"fmt"
	"time"
)

func TestContext() {
	basectx := context.Background()

	ctx, cancel := context.WithTimeout(basectx, 2*time.Second)

	defer cancel()

	for {
		fmt.Println("running in for...")
		time.Sleep(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			fmt.Println("context Done, because timeout is exceeded")
			return
		default:
		}
	}
}
