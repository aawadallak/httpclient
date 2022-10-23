package httpclient

import (
	"time"
)

type (
	CallOption func(*callConfig)

	callConfig struct {
		timeout time.Duration
	}
)
