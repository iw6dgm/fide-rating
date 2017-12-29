package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
)

const (
	PlayersFilename = `players_list_xml_foa.xml`
	PlayersDB       = `fide.db`
	DeleteSQL       = `DELETE FROM player`
	VacuumSQL       = `VACUUM`
	InsertSQL       = `INSERT INTO player (fideid,name,country,sex,title,w_title,o_title,foa_title,rating,games,k,rapid_rating,rapid_games,rapid_k,blitz_rating,blitz_games,blitz_k,birthday,flag) VALUES (? /*not nullable*/,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	SelectCountSQL  = `SELECT COUNT(*) FROM player`
)

type PlayersList struct {
	XMLName xml.Name `xml:"playerslist"`
	Players []Player `xml:"player"`
}

type Player struct {
	FideId      uint64 `xml:"fideid"`
	Name        string `xml:"name"`
	Country     string `xml:"country"`
	Sex         string `xml:"sex"`
	Title       string `xml:"title"`
	WTitle      string `xml:"w_title"`
	OTitle      string `xml:"o_title"`
	FoaTitle    string `xml:"foa_title"`
	Rating      uint   `xml:"rating"`
	Games       uint   `xml:"games"`
	K           uint8  `xml:"k"`
	RapidRating uint   `xml:"rapid_rating"`
	RapidGames  uint   `xml:"rapid_games"`
	RapidK      uint8  `xml:"rapid_k"`
	BlitzRating uint   `xml:"blitz_rating"`
	BlitzGames  uint   `xml:"blitz_games"`
	BlitzK      uint8  `xml:"blitz_k"`
	Birthday    uint16 `xml:"birthday"`
	Flag        string `xml:"flag"`
}

func main() {
	// Read XML file
	content := loadContent(PlayersFilename)

	pl := PlayersList{}
	// Parse XML content
	err := xml.Unmarshal(content, &pl)
	checkErr(err)

	// Open database connection
	db := dbOpen(PlayersDB)
	defer db.Close()

	// Clean up player table
	cleanUp(db)

	// Set up prepared statement to insert data
	stmt, _ := db.Prepare(InsertSQL)

	// Loop through decoded XML input data
	for _, p := range pl.Players {

		if p.Name == "" || p.Games == 0 {
			fmt.Printf("Skip player FIDE ID %d by having either Name or Games field empty\n", p.FideId)
			continue
		}

		stmt.Exec(
			// FIDE id
			p.FideId,
			// Basic info
			p.Name, p.Country, p.Sex,
			// Titles
			p.Title, p.WTitle, p.OTitle, p.FoaTitle,
			// Standard (classic) rating and K
			p.Rating, p.Games, p.K,
			// Rapid (classic) rating and K
			p.RapidRating, p.RapidGames, p.RapidK,
			// Blitz (classic) rating and K
			p.BlitzRating, p.BlitzGames, p.BlitzK,
			// Extra info
			p.Birthday, p.Flag)
	}
	var count int
	row := db.QueryRow(SelectCountSQL)
	row.Scan(&count)
	fmt.Printf("Total n. player(s) loaded : %d\n", count)
}

func cleanUp(db *sql.DB) {
	_, e1 := db.Exec(DeleteSQL)
	checkErr(e1)
	_, e2 := db.Exec(VacuumSQL)
	checkErr(e2)
}

func loadContent(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	checkErr(err)
	return content
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
