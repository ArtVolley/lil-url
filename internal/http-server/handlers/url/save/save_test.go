package save

import (
	"errors"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		name      string
		lilUrl    string
		url       string
		respError string
		mockError error
	}{
		{
			name:   "success",
			lilUrl: "test",
			url:    "https://test.ru",
		},
		{
			name: "generated lilUrl",
			url:  "",
		},
		{
			name:      "empty url",
			lilUrl:    "test",
			url:       "",
			respError: "field url is a required field",
		},
		{
			name:      "bad url",
			lilUrl:    "test",
			url:       "test",
			respError: "field url is not valid url",
		},
		{
			name:      "lil url already exists",
			lilUrl:    "test",
			url:       "https://test.ru",
			respError: "url already exists",
		},
		{
			name:      "failed to add",
			lilUrl:    "test",
			url:       "https://test.ru",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// urlSaverMock := mocks.NewUrlSaver(t)
		})
	}
}
