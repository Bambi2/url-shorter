package service

import (
	"fmt"
	"net/http"
)

func int64Pow(value int64, power int) int64 {
	if power == 0 {
		return 1
	}

	result := value
	for power > 1 {
		result *= value
		power--
	}

	return result
}

func urlLengthCheck(encodedUrl string, length int) error {
	if len(encodedUrl) != length {
		var msg string
		if len(encodedUrl) == 0 {
			msg = "empty request"
		} else {
			msg = fmt.Sprintf("invalid url length: length must be %d", length)
		}
		return &ServiceError{
			Msg:        msg,
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}
