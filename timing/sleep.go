package timing

import (
	"errors"
	"time"
)

var (
	timeSleep = time.Sleep
)

var (
	ErrExceededMaxTries = errors.New("timing: exceeded max tries")
)

// PollUntil receives a func that does some logic and it's gonna call this func every second until that func returns true.
// It calls this func based on maxTries param.
func PollUntil(fn func() bool, maxTries int) error {
	for i := 0; i < maxTries; i++ {
		timeSleep(1 * time.Minute)
		if fn() {
			return nil
		}
	}
	return ErrExceededMaxTries
}

type Poller struct {
	sleep func(time.Duration)
}

func (p *Poller) Until(fn func() bool, maxTries int) error {
	for i := 0; i < maxTries; i++ {
		p.sleep(1 * time.Minute)
		if fn() {
			return nil
		}
	}
	return ErrExceededMaxTries
}
