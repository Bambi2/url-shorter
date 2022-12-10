package service

import (
	"net/http"
	"strings"
)

const (
	base63Alphabet             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	base63AlphabetLength       = int64(len(base63Alphabet))
	base63EncodedLength        = 10
	Base63TenMaxId       int64 = 984930291881790848 //63^10-1
)

func base63Encode(id int64) string {
	i := base63EncodedLength - 1
	encoded := make([]byte, base63EncodedLength)

	for id > 0 {
		encoded[i] = base63Alphabet[id%base63AlphabetLength]
		id /= base63AlphabetLength
		i--
	}

	for i >= 0 {
		encoded[i] = 'a'
		i--
	}

	return string(encoded)
}

func base63Decode(encodedUrl string) (int64, error) {
	var id int64
	for i, symbol := range encodedUrl {
		pos := strings.IndexRune(base63Alphabet, symbol)
		if pos == -1 {
			return -1, &ServiceError{
				Msg:        "invalid character: " + string(symbol),
				StatusCode: http.StatusBadRequest,
			}
		}

		id += int64(pos) * int64Pow(base63AlphabetLength, base63EncodedLength-i-1)
	}

	return id, nil
}
