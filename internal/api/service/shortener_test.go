package service

import (
	"dh-url-shortener/internal/api/model"
	"errors"
	"testing"

	mocks "dh-url-shortener/.mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const longURL = "https://www.yemeksepeti.com/istanbul"

// TestShortener_Shorten should return error when url is empty
func TestShortener_Shorten_ShouldReturnErrorWhenLongURLIsEmpty(t *testing.T) {
	s := Shortener{}
	shortURL, err := s.Shorten("")

	assert.Error(t, err)
	assert.Equal(t, "", shortURL)
}

// TestShortener_Shorten should return first 7 chars of the hash
func TestShortener_Shorten_ShouldReturnFirstSevenCharsOfMd5Hash(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	s := Shortener{DB: mockDB}
	longURL := "https://www.yemeksepeti.com/istanbul"
	shortURL, err := s.Shorten(longURL)
	expected := "/05bf184"
	assert.Nil(t, err)
	assert.Equal(t, expected, shortURL)
}

// TestShortener_Shorten should return short url when if it is not used before
func TestShortener_Shorten_ShouldReturnShortUrlWhenIfItIsNotUsedBefore(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	s := Shortener{DB: mockDB}
	longURL := "https://www.yemeksepeti.com/istanbul"
	shortURL, err := s.Shorten(longURL)
	expected := "/05bf184"
	assert.Nil(t, err)
	assert.Equal(t, expected, shortURL)
}

// TestShortener_Shorten should return different short url when if short url already is used before
func TestShortener_Shorten_ShouldReturnDifferentShortUrlWhenIfItIsUsedBefore(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	gomock.InOrder(
		mockDB.EXPECT().Set("05bf184", model.RedirectionData{OriginalURL: longURL}).Return(errors.New("hash already exists")).Times(1),
		mockDB.EXPECT().Set("8d505df", model.RedirectionData{OriginalURL: longURL}).Return(nil).Times(1),
	)

	s := Shortener{DB: mockDB}
	shortURL, err := s.Shorten(longURL)
	expected := "/8d505df"
	assert.Nil(t, err)
	assert.Equal(t, expected, shortURL)
}

func TestShortener_Expand_ShouldReturnErrorWhenHashNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Get(gomock.Any()).Return(model.RedirectionData{}, errors.New("hash not found")).Times(1)

	s := Shortener{DB: mockDB}
	hash := "05bf184"
	originalURL, err := s.Expand(hash)

	assert.Error(t, err)
	assert.Equal(t, "", originalURL)
}

func TestShortener_Expand_ShouldReturnLongURL(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Get(gomock.Any()).Return(model.RedirectionData{OriginalURL: longURL}, nil).Times(1)
	mockDB.EXPECT().Hit(gomock.Any()).Return(nil).Times(1)

	s := Shortener{DB: mockDB}
	hash := "05bf184"
	url, err := s.Expand(hash)

	assert.Nil(t, err)
	assert.Equal(t, longURL, url)
}

func TestShortener_Expand_ShouldReturnErrorWhenCantIncreaseHit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Get(gomock.Any()).Return(model.RedirectionData{OriginalURL: longURL}, nil).Times(1)
	mockDB.EXPECT().Hit(gomock.Any()).Return(errors.New("key not found")).Times(1)

	s := Shortener{DB: mockDB}
	hash := "05bf184"
	url, err := s.Expand(hash)

	assert.Error(t, err)
	assert.Equal(t, "", url)
}

func TestShortener_List(t *testing.T) {
	data := map[string]model.RedirectionData{
		"05bf184": {OriginalURL: longURL, Hits: 5},
	}

	expectedResult := []model.ListData{
		{OriginalURL: longURL, Hits: 5, Hash: "05bf184"},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Data().Return(data).Times(1)

	s := Shortener{DB: mockDB}
	actualResult := s.List()
	assert.Equal(t, expectedResult, actualResult)
}

func TestShortener_List_ShouldReturnJSONArrayWhenDbIsEmpty(t *testing.T) {
	data := map[string]model.RedirectionData{}

	expectedResult := make([]model.ListData, 0, 1)

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDB := mocks.NewMockDB(controller)
	mockDB.EXPECT().Data().Return(data).Times(1)

	s := Shortener{DB: mockDB}
	actualResult := s.List()
	assert.Equal(t, expectedResult, actualResult)
}
