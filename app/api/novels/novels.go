package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

type (
	sentence struct {
		FirstLine  string `json:"first_line" db:"first_line"`
		SecondLine string `json:"second_line" db:"second_line"`
		Revision   string `json:"revision" db:"revision"`
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

func main() {
	e := echo.New()
	e.GET("/novels/:id", func(c echo.Context) (err error) {
		id := c.Param("id")

		// create a connection (e.g. "postgres", "mysql", or "sqlite3")
		connection, _ := dbr.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":3306)/"+os.Getenv("DATABASE"), nil)

		// create a session for each business unit of execution (e.g. a web request or goworkers job)
		session := connection.NewSession(nil)

		// get a record
		var novel novel
		var sentence []sentence
		novel.Sentence = sentence
		andCondition := dbr.And(dbr.Eq("novel.id", id))

		count, err := session.Select("novel.id", "novel.title", "novel.genre_id", "novel.summary", "novel.relay_limit", "novel.novelist_id", "novel.first_edition", "novelist.name").From("novel").Join("novelist", "novel.novelist_id = novelist.id").Where(andCondition).Load(&novel)

		if err != nil {
			panic(err.Error())
		}
		if count == 0 {
			errorResponse := response{Status: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound), Data: novel}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		session.Select("sentence.first_line", "sentence.second_line", "sentence.revision").From("sentence").Where("sentence.novel_id = ?", id).OrderBy("sentence.revision").Load(&novel.Sentence)

		response := response{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: novel}

		// json-ready, with dbr.null* types serialized like you want
		return c.JSON(http.StatusOK, response)
	})
	e.Logger.Fatal(e.Start(":80"))
}
