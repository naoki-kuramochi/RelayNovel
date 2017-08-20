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
		Revision   string `json:"-" db:"revision"`
		NovelID    int    `json:"-" db:"novel_id"`
		NovelistID int    `json:"-" db:"novelist_id"`
	}

	novel struct {
		ID           uint16     `json:"id" db:"id"`       // not null
		Title        string     `json:"title" db:"title"` // not null
		GenreID      uint16     `json:"genre_id" db:"genre_id"`
		Summary      string     `json:"summary" db:"summary"`             // not null
		RelayLimit   uint16     `json:"relay_limit" db:"relay_limit"`     // not null
		NovelistID   uint16     `json:"novelist_id" db:"novelist_id"`     // not null
		FirstEdition string     `json:"first_edition" db:"first_edition"` // not null
		NovelistName string     `json:"novelist" db:"name"`               // not null
		Sentence     []sentence `json:"sentence"`
	}

	response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    novel  `json:"data"`
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
	e.Logger.Fatal(e.Start(":80"))
}

// getNovels novelの詳細を返却
func getNovels(c echo.Context) (err error) {
	id := c.Param("id")

	novel := novel{Sentence: []sentence{}}

	session, err := createSession()
	if err != nil {
		log.Printf("alert: " + err.Error())
		return c.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError))
	}

	andCondition := dbr.And(dbr.Eq("novel.id", id))
	count, err := session.Select("novel.id", "novel.title", "novel.genre_id", "novel.summary", "novel.relay_limit", "novel.novelist_id", "novel.first_edition", "novelist.name").
		From("novel").
		Join("novelist", "novel.novelist_id = novelist.id").
		Where(andCondition).
		Load(&novel)
	if err != nil {
		log.Printf("alert: " + err.Error())
		return c.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError))
	}
	if count == 0 {
		return c.JSON(http.StatusNotFound, errorResponse(http.StatusNotFound))
	}

	sentenceApiResponse, err := fetchSentence(id)
	if err != nil {
		log.Printf("info: " + err.Error())
		return c.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest))
	}

	novel.Sentence = sentenceApiResponse.Data
	response := response{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: novel}
	return c.JSON(http.StatusOK, response)
}

// fetchSentence sentencesAPIを実行してdataを取得
func fetchSentence(id string) (sentenceApiResponse, error) {
	sentenceApiResponse := sentenceApiResponse{Data: []sentence{}}
	response, err := http.Get("http://sentences-api/novels/sentences/" + id)
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

// errorResponse エラー時のレスポンスを生成
func errorResponse(statusCode int) response {
	novel := novel{}
	novel.Sentence = []sentence{}
	response := response{Status: statusCode, Message: http.StatusText(statusCode), Data: novel}
	return response
}

// createSession dbr.Sessionを生成
func createSession() (*dbr.Session, error) {
	connection, err := dbr.Open("mysql", os.Getenv("RELAYNOVEL_SESSION"), nil)
	session := connection.NewSession(nil)
	return session, err
}
