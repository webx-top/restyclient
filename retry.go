package restyclient

import (
	"time"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/com"
)

func RetryAfter(c *resty.Client, response *resty.Response) (time.Duration, error) {
	retryAfter := response.Header().Get(`Retry-After`)
	if len(retryAfter) > 0 {
		if com.StrIsNumeric(retryAfter) {
			delaySeconds := com.Int64(retryAfter)
			return time.Second * time.Duration(delaySeconds), nil
		}
		t, err := time.Parse(time.RFC1123, retryAfter)
		if err != nil {
			return 0, err
		}
		return time.Until(t.Local()), nil
	}

	remaing := response.Header().Get(`X-Ratelimit-Remaining`)
	if len(remaing) > 0 {
		if com.Int64(remaing) > 0 {
			return 0, nil
		}
	}
	retryAfter = response.Header().Get(`X-Ratelimit-Reset`)
	if len(retryAfter) > 0 {
		if com.StrIsNumeric(retryAfter) {
			timestamp := com.Int64(retryAfter)
			return time.Until(time.Unix(timestamp, 0)), nil
		}
		t, err := time.Parse(time.RFC1123, retryAfter)
		if err != nil {
			return 0, err
		}
		return time.Until(t.Local()), nil
	}
	return 0, nil
}
