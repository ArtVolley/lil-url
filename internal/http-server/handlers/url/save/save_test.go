package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"lil-url/internal/http-server/handlers/url/save"
	"lil-url/internal/http-server/handlers/url/save/mocks"
	slogdiscard "lil-url/internal/lib/logger/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
			url:  "https://test.ru",
		},
		{
			name:      "empty url",
			lilUrl:    "test",
			url:       "",
			respError: "field Url is a required field",
		},
		{
			name:      "bad url",
			lilUrl:    "test",
			url:       "test",
			respError: "field Url is not valid url",
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

			urlSaverMock := mocks.NewUrlSaver(t)

			if tt.respError == "" || tt.mockError != nil {
				urlSaverMock.On(
					"SaveUrl",
					tt.url,
					mock.AnythingOfType("string")).
					Return(int64(1),
						tt.mockError).Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "lilUrl": "%s"}`, tt.url, tt.lilUrl)

			req, err := http.NewRequest(http.MethodPost,
				"/save",
				bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
