package handler

import (
	"bytes"
	"dh-url-shortener/internal/api/model"
	"dh-url-shortener/internal/api/service"
	"dh-url-shortener/internal/platform/db"
	"dh-url-shortener/internal/platform/snapshot"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader([]byte(`invalid json`)))

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
	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader([]byte(`{"url": "invalid url"}`)))

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
	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL)))) // nolint:gocritic

	handler.Shorten(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestShortenerHandler_Shorten_ShortenedURL(t *testing.T) {
	expectedShortenedURL := fmt.Sprintf(`{"url":"%s/tTeEsT"}`, shortURLDomain)
	shortenedURL := shortURLDomain + "/tTeEsT"
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().Shorten(gomock.Any()).Return(shortenedURL, nil).Times(1)

	handler := URLHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL)))) // nolint:gocritic

	handler.Shorten(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, expectedShortenedURL, resp.Body.String())
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
	req := httptest.NewRequest(http.MethodPost, "/short", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL)))) // nolint:gocritic
	handler.Shorten(resp, req)
	redirectionData, _ := InMemoryDB.Get("05bf184")
	expectedShortenedURL := fmt.Sprintf(`{"url":"%s/05bf184"}`, shortURLDomain)
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, expectedShortenedURL, resp.Body.String())
	assert.Equal(t, longURL, redirectionData.OriginalURL)
}

func TestUrlHandler_Expand_ShouldReturnBadRequestWhenHashSmallerThanSevenChar(t *testing.T) {
	handler := URLHandler{}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/05bf184", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/05bf184", nil)
	handler.Expand(resp, req)

	assert.Equal(t, http.StatusFound, resp.Code)
}

func TestURLHandler_List(t *testing.T) {
	testData := []model.ListData{
		{
			Hash:        "05bf184",
			OriginalURL: longURL,
			Hits:        4,
		},
	}
	expectedResp, _ := json.Marshal(testData)

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mocks.NewMockShortenerService(controller)
	mockShortenerService.EXPECT().List().Return(testData).Times(1)
	handler := URLHandler{ShortenerService: mockShortenerService}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	handler.List(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, string(expectedResp), resp.Body.String())
}

func BenchmarkURLHandler_Expand(b *testing.B) {
	inMemoryDB := db.NewInMemoryDB()
	shortenerService := service.Shortener{DB: inMemoryDB, ShortURLDomain: "http://localhost:8080"}
	h := URLHandler{ShortenerService: shortenerService}
	_ = inMemoryDB.Set("05bf184", model.RedirectionData{OriginalURL: longURL, Hits: 0})
	ss := snapshot.NewSnapshot("../db/test_snapshot.db", time.Second*5)
	_ = ss.Restore(inMemoryDB)
	go ss.SavePeriodically(inMemoryDB, nil)
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(16000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/05bf184", http.NoBody)
			h.Expand(resp, req)
		}
	})
}
