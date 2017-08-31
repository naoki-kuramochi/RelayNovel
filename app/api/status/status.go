package main

import (
	"github.com/comail/colog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"log"
	"net/http"
	"strconv"
)

type (
	response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    []int  `json:"data"`
	}
)

func init() {
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Flag: log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/:id", httpResponse)
	e.POST("/:id", httpResponse)
	e.GET("/status/:id", httpResponse)
	e.POST("/status/:id", httpResponse)

	e.Logger.Fatal(e.Start(":80"))
}

// httpResponse エラー時のレスポンスを返す
func httpResponse(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		id = http.StatusNotFound
	}

	// 存在しないstatus_codeが来た場合は404
	message := http.StatusText(id)
	if message == "" {
		id = http.StatusNotFound
		message = http.StatusText(http.StatusNotFound)
	}

	response := response{Status: id, Message: message, Data: make([]int, 0)} // make([]int, 0)で空の配列を生成
	return c.JSON(id, response)
}