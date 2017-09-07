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
	sentenceApiResponse struct {
		Status  int        `json:"status"`
		Message string     `json:"message"`
		Data    []sentence `json:"data"`
	}

	sentence struct {
		FirstLine  string `json:"first_line" db:"first_line"`
		SecondLine string `json:"second_line" db:"second_line"`
		Revision   string `json:"revision" db:"revision"`
		NovelID    int    `json:"-" db:"novel_id"`
		NovelistID int    `json:"-" db:"novelist_id"`
	}

	novel struct {
		ID           int        `json:"id" db:"id"`       // not null
		Title        string     `json:"title" db:"title"` // not null
		GenreID      int        `json:"genre_id" db:"genre_id"`
		Summary      string     `json:"summary" db:"summary"`             // not null
		RelayLimit   int        `json:"relay_limit" db:"relay_limit"`     // not null
		NovelistID   int        `json:"novelist_id" db:"novelist_id"`     // not null
		FirstEdition string     `json:"first_edition" db:"first_edition"` // not null
		NovelistName string     `json:"novelist" db:"name"`               // not null
		Sentence     []sentence `json:"sentence"`
	}

	response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    novel  `json:"data"`
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

	e.GET("/novels/:id", getNovels)
	//e.POST("/novels/:id", getNovels)
	e.Logger.Fatal(e.Start(":80"))
}

// getNovels novelの詳細を返却
func getNovels(c echo.Context) (err error) {
	id := c.Param("id")

	novel := novel{Sentence: []sentence{}}

	session, err := createSession()
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}

	andCondition := dbr.And(dbr.Eq("novel.id", id))
	count, err := session.Select("novel.id", "novel.title", "novel.genre_id", "novel.summary", "novel.relay_limit", "novel.novelist_id", "novel.first_edition", "novelist.name").
		From("novel").
		Join("novelist", "novel.novelist_id = novelist.id").
		Where(andCondition).
		Load(&novel)
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}
	if count == 0 {
		errors, _ := fetchErrorResponse(http.StatusNotFound)
		return c.JSON(http.StatusNotFound, errors)
	}

	sentenceApiResponse, err := fetchSentence(id)
	if err != nil {
		log.Printf("info: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusBadRequest)
		return c.JSON(http.StatusBadRequest, errors)
	}

	novel.Sentence = sentenceApiResponse.Data
	response := response{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: novel}
	return c.JSON(http.StatusOK, response)
}

// fetchSentence sentencesAPIを実行してdataを取得
func fetchSentence(id string) (sentenceApiResponse, error) {
	sentenceApiResponse := sentenceApiResponse{Data: []sentence{}}
	response, err := http.Get("http://sentences-api/novels/" + id + "/sentences")
	defer response.Body.Close()
	if err != nil {
		return sentenceApiResponse, err
	}

	err = json.NewDecoder(response.Body).Decode(&sentenceApiResponse)
	if err != nil {
		return sentenceApiResponse, err
	}

	return sentenceApiResponse, nil
}

// fetchErrorResponse エラー時のレスポンスを生成
func fetchErrorResponse(statusCode int) (errors, error) {
	errorsApiResponse := errors{}
	response, err := http.Get("http://status-api/" + strconv.Itoa(statusCode))
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
