package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// http.statusが200の場合
	getSentence = `{"status":200,"message":"OK","data":[]}`
)

// TestGetSentence 文章を取得
func TestGetSentence(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/novels/:novelID/sentences")
	c.SetParamNames("NovelID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, httpResponse(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, okJSON, rec.Body.String())
	}
}

// TestFetchErrorResponse エラーレスポンスを返すテスト
func TestFetchErrorResponse(t *testing.T) {
	statusCode := 500
	response, err := fetchErrorResponse(statusCode)
	assert.Equal(t, response.Status, statusCode)
	assert.Nil(t, err)
}

// TestCreateSessio DB接続のテスト
func TestCreateSession(t *testing.T) {
	_, err := createSession()
	// Assertions
	assert.Nil(t, err)
}
