package poller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Response reddit.Response

func (r Response) RequestsUsed() (int, error) {
	usedStr := r.Header.Get("X-Ratelimit-Used")
	if usedStr == "" {
		return 0, errors.New("header does not contain X-Ratelimit-Used")
	}

	n, err := strconv.Atoi(usedStr)
	if err != nil {
		return 0, fmt.Errorf("could not convert: %s to integer", usedStr)
	}
	return n, nil
}

func (r Response) RequestsRemaining() (int, error) {
	usedStr := r.Header.Get("X-Ratelimit-Remaining")
	if usedStr == "" {
		return 0, errors.New("header does not contain X-Ratelimit-Remaining")
	}

	f, err := strconv.ParseFloat(usedStr, 64)
	if err != nil {
		return 0, fmt.Errorf("could not convert: %s to float64", usedStr)
	}
	return int(f), nil
}

func (r Response) SecondsUntilReset() (int, error) {
	usedStr := r.Header.Get("X-Ratelimit-Reset")
	if usedStr == "" {
		return 0, errors.New("header does not contain X-Ratelimit-Reset")
	}

	n, err := strconv.Atoi(usedStr)
	if err != nil {
		return 0, fmt.Errorf("could not convert: %s to integer", usedStr)
	}
	return n, nil
}
