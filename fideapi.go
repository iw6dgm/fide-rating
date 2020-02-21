package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Player struct {
	FideId      uint64 `json:"fideid"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Sex         string `json:"sex"`
	Title       string `json:"title"`
	WTitle      string `json:"w_title"`
	OTitle      string `json:"o_title"`
	FoaTitle    string `json:"foa_title"`
	Rating      uint   `json:"rating"`
	Games       uint   `json:"games"`
	K           uint8  `json:"k"`
	RapidRating uint   `json:"rapid_rating"`
	RapidGames  uint   `json:"rapid_games"`
	RapidK      uint8  `json:"rapid_k"`
	BlitzRating uint   `json:"blitz_rating"`
	BlitzGames  uint   `json:"blitz_games"`
	BlitzK      uint8  `json:"blitz_k"`
	Birthday    uint16 `json:"birthday"`
	Flag        string `json:"flag"`
}

const (
	PlayersDB = `fide.db`
)

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

//var db *sql.DB

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})

	r.GET("/players/:playerID", func(c *gin.Context) {
		var player *Player
		var err error

		if playerID := c.Param("playerID"); playerID != "" {
			if id, e := strconv.ParseUint(playerID, 10, 64); e == nil {
				player, err = dbGetPlayer(id)

				if err == sql.ErrNoRows {
					c.JSON(http.StatusNotFound, gin.H{"message": ErrNotFound})
					return
				}

			} else {
				c.JSON(http.StatusNotFound, gin.H{"message": ErrNotFound})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": ErrNotFound})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"player": player,
		})
	})

	http.ListenAndServe(":7373", r)
}

func dbGetPlayer(playerID uint64) (*Player, error) {
	// Open database connection
	db := dbOpen(PlayersDB)
	defer db.Close()

	p := Player{FideId: playerID}
	err := db.QueryRow("SELECT name,country,sex,title,w_title,o_title,foa_title,rating,games,k,rapid_rating,rapid_games,rapid_k,blitz_rating,blitz_games,blitz_k,birthday,flag FROM player WHERE fideid=?", playerID).Scan(&p.Name, &p.Country, &p.Sex, &p.Title, &p.WTitle, &p.OTitle, &p.FoaTitle, &p.Rating, &p.Games, &p.K, &p.RapidRating, &p.RapidGames, &p.RapidK, &p.BlitzRating, &p.BlitzGames, &p.BlitzK, &p.Birthday, &p.Flag)
	if err != nil {
		return &p, err
	}
	return &p, nil
}

func dbOpen(conn string) *sql.DB {
	db, err := sql.Open("sqlite3", conn)
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
