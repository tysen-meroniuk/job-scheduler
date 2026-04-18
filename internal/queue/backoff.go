package queue

import (
	"math"
	"math/rand/v2"
	"time"
)

// ExpBackoff returns the delay before the next retry of a job that has been
// attempted `attempts` times. Formula: min(2^attempts, 300) seconds + 0..5s jitter.
// The cap at 300s prevents multi-hour delays on long-failing jobs; jitter
// prevents synchronized thundering-herd retries across workers.
func ExpBackoff(attempts int) time.Duration {
	base := math.Min(math.Pow(2, float64(attempts)), 300)
	jitter := rand.Float64() * 5
	return time.Duration((base + jitter) * float64(time.Second))
}
