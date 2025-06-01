package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Isbn     string    `json:"isbn"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	fmt.Println("ID:", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	for _, movie := range movies {

		if movie.ID == id {

			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func addMovies(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = rand.Intn(100)
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	MovieID := mux.Vars(r)["id"]
	var updateMovie Movie
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&updateMovie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(MovieID)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	for i, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:i], movies[i+1:]...)
			updateMovie.ID = movie.ID
			movies = append(movies, updateMovie)
			json.NewEncoder(w).Encode(updateMovie)
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	MovieID := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	for i, movie := range movies {
		if strconv.Itoa(movie.ID) == MovieID {
			movies = append(movies[:i], movies[i+1:]...)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Movie not found",
		"status":  "error",
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello, World!",
		"status":  "success",
	})
}
func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: rand.Intn(100), Title: "Kushal The Geams Bond", Isbn: "438743", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: rand.Intn(100), Title: " Kushal The ", Isbn: "438743", Director: &Director{Firstname: "Jane", Lastname: "Doe"}})

	r.HandleFunc("/", hello).Methods("GET")
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", addMovies).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
