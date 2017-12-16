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
	novelist struct {
		ID   string `json:"id" db:"id"`     // not null
		Name string `json:"name" db:"name"` // not null
	}

	response struct {
		Status  int        `json:"status"`
		Message string     `json:"message"`
		Data    []novelist `json:"data"`
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

	e.GET("/novelists", getNovelists)
	//e.GET("/novelists/:id", getNovelistsDetail)
	//e.POST("/novels/:id", getNovels)
	e.Logger.Fatal(e.Start(":80"))
}

// getNovels novelの詳細を返却
func getNovelists(c echo.Context) (err error) {
	novelist := []novelist{}

	session, err := createSession()
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}

	count, err := session.Select("novelist.id", "novelist.name", "novelist.create", "novelist.update").
		From("novelist").
		Load(&novelist)
	if err != nil {
		log.Printf("alert: " + err.Error())
		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, errors)
	}
	if count == 0 {
		errors, _ := fetchErrorResponse(http.StatusNotFound)
		return c.JSON(http.StatusNotFound, errors)
	}

	response := response{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: novelist}
	return c.JSON(http.StatusOK, response)
}

// getNovels novelの詳細を返却
//func getNovelistsDetail(c echo.Context) (err error) {
//	id := c.Param("id")
//
//	novel := novel{Sentence: []sentence{}}
//
//	session, err := createSession()
//	if err != nil {
//		log.Printf("alert: " + err.Error())
//		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
//		return c.JSON(http.StatusInternalServerError, errors)
//	}
//
//	andCondition := dbr.And(dbr.Eq("novel.id", id))
//	count, err := session.Select("novel.id", "novel.title", "novel.genre_id", "novel.summary", "novel.relay_limit", "novel.novelist_id", "novel.first_edition", "novelist.name").
//		From("novel").
//		Join("novelist", "novel.novelist_id = novelist.id").
//		Where(andCondition).
//		Load(&novel)
//	if err != nil {
//		log.Printf("alert: " + err.Error())
//		errors, _ := fetchErrorResponse(http.StatusInternalServerError)
//		return c.JSON(http.StatusInternalServerError, errors)
//	}
//	if count == 0 {
//		errors, _ := fetchErrorResponse(http.StatusNotFound)
//		return c.JSON(http.StatusNotFound, errors)
//	}
//
//	response := response{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: novel}
//	return c.JSON(http.StatusOK, response)
//}

// fetchErrorResponse エラー時のレスポンスを生成
func fetchErrorResponse(statusCode int) (errors, error) {
	errorsAPIResponse := errors{}
	response, err := http.Get("http://status-api/" + strconv.Itoa(statusCode))
	defer response.Body.Close()
	if err != nil {
		return errorsAPIResponse, err
	}

	err = json.NewDecoder(response.Body).Decode(&errorsAPIResponse)
	if err != nil {
		return errorsAPIResponse, err
	}

	return errorsAPIResponse, nil
}

// createSession dbr.Sessionを生成
func createSession() (*dbr.Session, error) {
	connection, err := dbr.Open("mysql", os.Getenv("RELAYNOVEL_SESSION"), nil)
	session := connection.NewSession(nil)
	return session, err
}
