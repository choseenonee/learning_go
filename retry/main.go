package retry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

var ans string = "5"

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {

		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				ans = "cancelled!!!"
				fmt.Printf("aaa %v", ans)
				return "", ctx.Err()
			}
		}
	}
}

var count int

func EmulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}

func TestRetry() {
	r := Retry(EmulateTransientError, 5, 2*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1999)
	defer cancel()
	go r(ctx)
	select {
	case <-time.After(time.Second * 3):
		//cancel()
		fmt.Printf("%v", ans)
	}
}
