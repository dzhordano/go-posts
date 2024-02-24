package utils

import (
	"time"
)

func DoWithTries(fn func() error, atts int, delay time.Duration) (err error) {
	for atts < 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			atts--

			continue
		}

		return nil
	}

	return
}
