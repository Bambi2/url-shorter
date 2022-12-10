package service

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/bambi2/url-shorter/internal/repository"
	mock_repository "github.com/bambi2/url-shorter/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestService_encoderbase63(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockEncoder, url string)

	testTable := []struct {
		name               string
		inputUrl           string
		mockBehavior       mockBehavior
		expectedEncodedURL string
		expectedErrorMsg   string
		expectedStatusCode int
	}{
		{
			name:               "invalidUrl",
			inputUrl:           "sss",
			mockBehavior:       func(s *mock_repository.MockEncoder, url string) {},
			expectedEncodedURL: "",
			expectedErrorMsg:   "invalid url: sss",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "urlExists",
			inputUrl: "https://www.google.com",
			mockBehavior: func(s *mock_repository.MockEncoder, url string) {
				s.EXPECT().IfExistsBase63(url).Return(int64(50), nil)
			},
			expectedEncodedURL: base63Encode(50),
			expectedErrorMsg:   "",
			expectedStatusCode: -1,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			encoder := mock_repository.NewMockEncoder(c)
			testCase.mockBehavior(encoder, testCase.inputUrl)

			repo := &repository.Repository{
				Encoder: encoder,
			}

			service := NewService(repo)

			//Testing
			encodedUrl, err := service.Encoder.Base63(testCase.inputUrl)
			var serviceErr *ServiceError
			if !errors.As(err, &serviceErr) {
				serviceErr = &ServiceError{Msg: "", StatusCode: -1}
			}

			assert.Equal(t, encodedUrl, testCase.expectedEncodedURL)
			assert.Equal(t, serviceErr.Msg, testCase.expectedErrorMsg)
			assert.Equal(t, serviceErr.StatusCode, testCase.expectedStatusCode)
		})
	}
}

func TestService_decoderbase63(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockDecoder, encodedUrl string)

	testTable := []struct {
		name                string
		inputEncodedUrl     string
		mockBehavior        mockBehavior
		expectedOriginalURL string
		expectedErrorMsg    string
		expectedStatusCode  int
	}{
		{
			name:                "wrongLength",
			inputEncodedUrl:     "123456789",
			mockBehavior:        func(s *mock_repository.MockDecoder, encodedUrl string) {},
			expectedOriginalURL: "",
			expectedErrorMsg:    fmt.Sprintf("invalid url length: length must be %d", base63EncodedLength),
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name:                "zeroLength",
			inputEncodedUrl:     "",
			mockBehavior:        func(s *mock_repository.MockDecoder, encodedUrl string) {},
			expectedOriginalURL: "",
			expectedErrorMsg:    "empty request",
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name:            "noSuchUrl",
			inputEncodedUrl: "freeeMoney",
			mockBehavior: func(s *mock_repository.MockDecoder, encodedUrl string) {
				id, _ := base63Decode("freeeMoney")
				s.EXPECT().GetBase63(id).Return("", repository.ErrNoSuchURL)
			},
			expectedOriginalURL: "",
			expectedErrorMsg:    "no such url",
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name:                "invalidCharacter",
			inputEncodedUrl:     "123456789+",
			mockBehavior:        func(s *mock_repository.MockDecoder, encodedUrl string) {},
			expectedOriginalURL: "",
			expectedErrorMsg:    "invalid character: +",
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name:            "decoding",
			inputEncodedUrl: "asd_asdasd",
			mockBehavior: func(s *mock_repository.MockDecoder, encodedUrl string) {
				id, _ := base63Decode("asd_asdasd")
				s.EXPECT().GetBase63(id).Return("https://www.google.com", nil)
			},
			expectedOriginalURL: "https://www.google.com",
			expectedErrorMsg:    "",
			expectedStatusCode:  -1,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			decoder := mock_repository.NewMockDecoder(c)
			testCase.mockBehavior(decoder, testCase.inputEncodedUrl)

			repo := &repository.Repository{
				Decoder: decoder,
			}

			service := NewService(repo)

			//Testing
			originalUrl, err := service.Decoder.Base63(testCase.inputEncodedUrl)
			var serviceErr *ServiceError
			if !errors.As(err, &serviceErr) {
				serviceErr = &ServiceError{Msg: "", StatusCode: -1}
			}

			assert.Equal(t, originalUrl, testCase.expectedOriginalURL)
			assert.Equal(t, serviceErr.Msg, testCase.expectedErrorMsg)
			assert.Equal(t, serviceErr.StatusCode, testCase.expectedStatusCode)
		})
	}
}
