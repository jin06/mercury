package clients

import "time"

func DefaultOptions() *Options {
	return &Options{
		PublishTimeout:  time.Second * 2,
		MaxPublishTimes: 5,
	}
}

type Options struct {
	PublishTimeout  time.Duration
	MaxPublishTimes int
}
