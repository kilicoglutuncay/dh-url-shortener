package handler

import (
	"bytes"
	"dh-url-shortener/db"
	"dh-url-shortener/service"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "dh-url-shortener/.mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	shortURLDomain = "http://localhost:8080"
	longURL        = "https://www.yemeksepeti.com/istanbul"
)

func TestShortenerHandler_Shorten_ShouldReturnBadRequestWhenShortenRequestIsNotContainsValidJSON(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return("", nil).Times(0)

	handler := URLHandler{
		ShortenerService: mockShortenerService,
	}
	resp := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(`invalid json`)))

	handler.Shorten(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestShortenerHandler_Shorten_ShouldReturnBadRequestWhenShortenRequestIsNotValid(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return("", nil).Times(0)

	handler := URLHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(`{"url": "invalid url"}`)))

	handler.Shorten(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestShortenerHandler_Shorten_ShouldReturnInternalServerErrorWhenShortenerServiceReturnsError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return("", errors.New("service error")).Times(1)

	handler := URLHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL))))

	handler.Shorten(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestShortenerHandler_Shorten_ShortenedURL(t *testing.T) {
	shortenedURL := shortURLDomain + "/tTeEsT"
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return(shortenedURL, nil).Times(1)

	handler := URLHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL))))

	handler.Shorten(resp, req)

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

// TestShortenerHandler_Create tests integration of short url creation process
func TestShortenerHandler_Create(t *testing.T) {
	InMemoryDB := db.NewInMemoryDB()
	svc := service.Shortener{DB: InMemoryDB, ShortURLDomain: shortURLDomain}
	handler := URLHandler{ShortenerService: svc}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/short", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL))))
	handler.Shorten(resp, req)
	redirectionData, _ := InMemoryDB.Get("05bf184")
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, shortURLDomain+"/05bf184", resp.Body.String())
	assert.Equal(t, longURL, redirectionData.OriginalURL)
}

func TestUrlHandler_Expand_ShouldReturnBadRequestWhenHashSmallerThanSevenChar(t *testing.T) {
	handler := URLHandler{}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	handler.Expand(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUrlHandler_Expand_ShouldReturnStatusNotFoundWhenServiceReturnsError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Expand(gomock.Any()).Return("", errors.New("hash not found")).Times(1)

	handler := URLHandler{ShortenerService: mockShortenerService}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/05bf184", nil)
	handler.Expand(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUrlHandler_Expand_ShouldReturnStatusFoundWhenServiceReturnsLongURL(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Expand(gomock.Any()).Return(longURL, nil).Times(1)

	handler := URLHandler{ShortenerService: mockShortenerService}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/05bf184", nil)
	handler.Expand(resp, req)

	assert.Equal(t, http.StatusFound, resp.Code)
}
