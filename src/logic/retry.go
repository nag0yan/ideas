package logic

import "time"

// Retry is a function that retries the given function up to maxRetries times with a delay between each retry.
// It returns the result of the function if it succeeds, or an error if it fails after all retries.
// The function takes a function that returns an error and a delay duration.
// The function is useful for handling transient errors that may succeed if retried.
// It is a simple implementation of the retry pattern.

type Timer interface {
	Sleep(d time.Duration)
}

func Retry[T any](fn func() (T, error), maxRetries int, delay time.Duration, timer Timer) (T, error) {
	var result T
	var err error

	for i := 0; i < maxRetries; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		timer.Sleep(min(delay<<i, time.Millisecond*100000)) // Exponential backoff
	}

	return result, err
}
