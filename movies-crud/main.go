package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`       // Field 'ID' will map to 'id' in JSON
	ISBN     string    `json:"isbn"`     // Field 'ISBN' will map to 'isbn' in JSON
	Title    string    `json:"title"`    // Field 'title' (lowercased) will map to 'title' in JSON
	Director *Director `json:"director"` // Field 'Director' will map to 'director' in JSON
}

type Director struct {
	firstName string `json""firstName"`
	lastName  string `json"lastName"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

func createNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.IntN(1000000))
	movies := append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: "1", ISBN: "23242", Title: "Movie One", Director: &Director{
		firstName: "John",
		lastName:  "Doe",
	}})

	movies = append(movies, Movie{
		ID:    "2",
		ISBN:  "2323",
		Title: "Movie Two",
		Director: &Director{
			firstName: "Mike",
			lastName:  "Tyson",
		},
	})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/[id]", getMovie).Methods("GET")
	r.HandleFunc("/movies", createNew).Methods("POST")
	r.HandleFunc("/movies/[id]", updateMovie).Methods("PUT")
	r.HandleFunc("movies/[id]", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))

}
