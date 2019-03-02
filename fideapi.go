package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"time"
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

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

//var db *sql.DB

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	r.Route("/players", func(r chi.Router) {
		r.Route("/{playerID}", func(r chi.Router) {
			r.Use(PlayerCtx)
			r.Get("/", getPlayer)
		})
	})

	http.ListenAndServe(":7373", r)
}

func PlayerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var player *Player
		var err error

		if playerID := chi.URLParam(r, "playerID"); playerID != "" {
			if id, e := strconv.ParseUint(playerID, 10, 64); e == nil {
				player, err = dbGetPlayer(id)

				if err == sql.ErrNoRows {
					render.Render(w, r, ErrNotFound)
					return
				}

			} else {
				render.Render(w, r, ErrNotFound)
				return
			}
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "player", player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (p *Player) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	player, ok := ctx.Value("player").(*Player)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	render.Render(w, r, player)
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
