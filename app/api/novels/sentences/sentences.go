package main

import (
	"github.com/comail/colog"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type (
	sentence struct {
		FirstLine  string `json:"first_line" db:"first_line"`
		SecondLine string `json:"second_line" db:"second_line"`
		Revision   string `json:"revision" db:"revision"`
		NovelID    int    `json:"novel_id" db:"novel_id"`
		NovelistID int    `json:"novelist_id" db:"novelist_id"`
	}

	response struct {
		Status  int        `json:"status"`
		Message string     `json:"message"`
		Data    []sentence `json:"data"`
	}

	errors struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Data    []interface{} `json:"data"`
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

	e.GET("/novels/sentences/:id", getSentence)
	e.POST("/novels/sentences", postSentence)

	e.Logger.Fatal(e.Start(":80"))
}

// getSentence novelに紐づく文章を取得
func getSentence(c echo.Context) (err error) {
	id := c.Param("id")
	sentence := []sentence{}

	session, err := createSession()
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}

	condition := dbr.And(dbr.Eq("sentence.novel_id", id))

	count, err := session.Select("sentence.novel_id", "sentence.novelist_id", "sentence.first_line", "sentence.second_line", "sentence.revision").
		From("sentence").
		Where(condition).
		OrderBy("sentence.revision").
		Load(&sentence)
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}
	if count == 0 {
		errors, _ := fetchErrorResponse(http.StatusNotFound)
		return c.JSON(http.StatusNotFound, errors)
	}

	response := response{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: sentence}
	return c.JSON(http.StatusOK, response)
}

// postSentence novelに紐づく文章を投稿
func postSentence(c echo.Context) (err error) {
	selectSentence := []sentence{}
	insertSentence := sentence{}

	err = c.Bind(&insertSentence)
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusBadRequest)
		return c.JSON(http.StatusBadRequest, errors)
	}

	session, err := createSession()
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}

	// transction
	transaction, err := session.Begin()
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}
	defer transaction.RollbackUnlessCommitted()

	result, err := transaction.InsertInto("sentence").
		Columns("novel_id", "novelist_id", "first_line", "second_line").
		Values(insertSentence.NovelID, insertSentence.NovelistID, insertSentence.FirstLine, insertSentence.SecondLine).
		Exec()
	// commit
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}
	transaction.Commit()

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}

	condition := dbr.And(dbr.Eq("sentence.id", lastInsertID))
	count, err := session.Select("sentence.novel_id", "sentence.novelist_id", "sentence.first_line", "sentence.second_line", "sentence.revision").
		From("sentence").
		Where(condition).
		OrderBy("sentence.revision DESC").
		Load(&selectSentence)
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}
	if count == 0 {
		errors, _ := fetchErrorResponse(http.StatusNotFound)
		return c.JSON(http.StatusNotFound, errors)
	}

	response := response{Status: http.StatusCreated, Message: http.StatusText(http.StatusCreated), Data: selectSentence}
	return c.JSON(http.StatusCreated, response)
}

// fetchErrorResponse エラー時のレスポンスを生成
func fetchErrorResponse(statusCode int) (errors, error) {
	errorsApiResponse := errors{}
	response, err := http.Get("http://errors-api/" + strconv.Itoa(statusCode))
	defer response.Body.Close()
	if err != nil {
		return errorsApiResponse, err
	}

	err = json.NewDecoder(response.Body).Decode(&errorsApiResponse)
	if err != nil {
		return errorsApiResponse, err
	}

	return errorsApiResponse, nil
}

// createSession dbr.Sessionを生成
func createSession() (*dbr.Session, error) {
	connection, err := dbr.Open("mysql", os.Getenv("RELAYNOVEL_SESSION"), nil)
	session := connection.NewSession(nil)
	return session, err
}
