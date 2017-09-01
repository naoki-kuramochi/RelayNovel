package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	// http.statusが200の場合
	okJSON = `{"status":200,"message":"OK","data":[]}`
	// http.statusが404もしくはXXX等の存在しない値の場合
	notFoundJSON = `{"status":404,"message":"Not Found","data":[]}`
	// http.statusが500の場合
	internalServerErrorJSON = `{"status":500,"message":"Internal Server Error","data":[]}`
)

// TestHttpResponseOk status_code=200でリクエストした時の返却値をテスト
func TestHttpResponseOk(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:statusCode")
	c.SetParamNames("statusCode")
	c.SetParamValues("200")

	// Assertions
	if assert.NoError(t, httpResponse(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, okJSON, rec.Body.String())
	}
}

// TestHttpResponseNotFound status_code=404でリクエストした時の返却値をテスト
func TestHttpResponseNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:statusCode")
	c.SetParamNames("statusCode")
	c.SetParamValues("404")

	// Assertions
	if assert.NoError(t, httpResponse(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundJSON, rec.Body.String())
	}
}

// TestHttpResponseInternalServerError status_code=500でリクエストした時の返却値をテスト
func TestHttpResponseInternalServerError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:statusCode")
	c.SetParamNames("statusCode")
	c.SetParamValues("500")

	// Assertions
	if assert.NoError(t, httpResponse(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, internalServerErrorJSON, rec.Body.String())
	}
}

// TestHttpResponseXXX status_code=XXXでリクエストした時の返却値をテスト
func TestHttpResponseXXX(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:statusCode")
	c.SetParamNames("statusCode")
	c.SetParamValues("XXX")

	// Assertions
	if assert.NoError(t, httpResponse(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundJSON, rec.Body.String())
	}
}
