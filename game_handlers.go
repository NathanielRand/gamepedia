package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Game struct holds game related data
type Game struct {
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	Rating      string `json:"rating"`
	Link        string `json:"link"`
}

var games []Game

func getGameHandler(w http.ResponseWriter, r *http.Request) {
	// Convert the "games" variable to json.
	gameListBytes, err := json.Marshal(games)

	// If there is an error, print it to console,
	// and return a server error response to the user.
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If all goes well, write the JSON list of games to the response.
	w.Write(gameListBytes)
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new instance of Game
	game := Game{}

	// We send all our data as HTML form data. 'ParseForm'
	// method of the request parses the form values.
	err := r.ParseForm()

	// In case of any error, we respond with an error to the error.
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the information about the game from the form info.
	game.Title = r.Form.Get("title")
	game.Genre = r.Form.Get("genre")
	game.Description = r.Form.Get("description")
	game.Rating = r.Form.Get("rating")
	game.Link = r.Form.Get("link")

	// Append our existing list of games with a new entry.
	games = append(games, game)

	// Finally, we redirect the user to the original HTML page
	// (located at '/assets/'), using the http libraries 'Redirect' method.
	http.Redirect(w, r, "/assets/", http.StatusFound)
}
