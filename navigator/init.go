package navigator

import (
	"time"

	"golang.org/x/time/rate"
)

// logLimiter controls the rate of log output.
var logLimiter *rate.Limiter

// Initialize the rate limiter with a limit of 1 event per second and a burst size of 5.
func init() {
	logLimiter = rate.NewLimiter(rate.Every(time.Second), 5)
}
