package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dh-url-shortener/.mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestShortenerHandler_Shorten_ShouldReturnBadRequestWhenShortenRequestIsNotContainsValidJSON(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return("", nil).Times(0)

	handler := ShortenerHandler{
		ShortenerService: mockShortenerService,
	}
	resp := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(`invalid json`)))

	handler.Create(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestShortenerHandler_Shorten_ShouldReturnBadRequestWhenShortenRequestIsNotValid(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return("", nil).Times(0)

	handler := ShortenerHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(`{"url": "invalid url"}`)))

	handler.Create(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestShortenerHandler_Shorten_ShouldReturnInternalServerErrorWhenShortenerServiceReturnsError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return("", errors.New("service error")).Times(1)

	handler := ShortenerHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(`{"url": "https://www.yemeksepeti.com/istanbul"}`)))

	handler.Create(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestShortenerHandler_Shorten_ShortenedURL(t *testing.T) {
	shortenedURL := "http://localhost:8080/tTeEsT"
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return(shortenedURL, nil).Times(1)

	handler := ShortenerHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(`{"url": "https://www.yemeksepeti.com/istanbul"}`)))

	handler.Create(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, shortenedURL, resp.Body.String())
}

func TestShortenRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		sr      ShortenRequest
		wantErr bool
	}{
		{
			name:    "long url is empty",
			sr:      ShortenRequest{URL: ""},
			wantErr: true,
		},
		{
			name:    "original url does not contain protocol",
			sr:      ShortenRequest{URL: "yemeksepeti.com"},
			wantErr: true,
		},
		{
			name:    "original url is a valid url",
			sr:      ShortenRequest{URL: "https://yemeksepeti.com"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualErr := tt.sr.validate()
			assert.Equal(t, tt.wantErr, actualErr != nil)
		})
	}
}
