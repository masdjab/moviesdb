package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"moviesdb.com/database"
	"moviesdb.com/model"
)

const (
	maxTitleLength       = 50
	maxDescriptionLength = 200
	maxImageUrlLength    = 200
)

// GET /ping
func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "pong")
	}
}

// POST /login
func LoginHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		uname := strings.Trim(params.Get("username"), " ")
		psswd := strings.Trim(params.Get("password"), " ")

		if uname == "" {
			badRequest(w, "Missing parameter: 'username'")
			return
		}

		if psswd == "" {
			badRequest(w, "Missing parameter: 'password'")
			return
		}

		user, err := db.LoginUser(uname, psswd)
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		success(w, user)
	}
}

// GET /logout
func LogoutHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getCurrentUser(db, r)
		if err == nil {
			db.LogoutUser(user.UserId)
		}
		success(w, nil)
	}
}

// GET /movies
func GetMoviesHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authRequest(db, r)
		if err != nil {
			notAuthorized(w, err.Error())
			return
		}

		params := r.URL.Query()
		keywords := params.Get("keywords")
		movies, err := db.GetMovies(keywords)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		success(w, movies)
	}
}

// GET /movies/{id}
func GetMovieDetailHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authRequest(db, r)
		if err != nil {
			notAuthorized(w, err.Error())
			return
		}

		movieId, err := getMovieIdFromRequest(r)
		if err != nil {
			badRequest(w, "Invalid ID")
			return
		}

		movie, err := db.GetMovieDetail(movieId)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
		}

		if movie == nil {
			badRequest(w, "Not found")
			return
		}

		success(w, *movie)
	}
}

// POST /movies
func PostMovieHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authRequest(db, r)
		if err != nil {
			notAuthorized(w, err.Error())
			return
		}

		movie, err := getMovieDataFromRequest(r)
		if err != nil {
			badRequest(w, err.Error())
			return
		}

		if err := validateMovieRequestData(movie); err != nil {
			badRequest(w, err.Error())
			return
		}

		id, err := db.InsertMovie(movie)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		savedMovie, err := db.GetMovieDetail(id)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		success(w, *savedMovie)
	}
}

// PATCH /movies/{id}
func PatchMovieHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authRequest(db, r)
		if err != nil {
			notAuthorized(w, err.Error())
			return
		}

		movieId, err := getMovieIdFromRequest(r)
		if err != nil {
			badRequest(w, "Invalid ID")
			return
		}

		movie, err := getMovieDataFromRequest(r)
		if err := validateMovieRequestData(movie); err != nil {
			badRequest(w, err.Error())
			return
		}

		err = db.UpdateMovie(movieId, movie)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		savedMovie, err := db.GetMovieDetail(movieId)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		success(w, *savedMovie)
	}
}

// DELETE /movies/{id}
func DeleteMovieHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authRequest(db, r)
		if err != nil {
			notAuthorized(w, err.Error())
			return
		}

		movieId, err := getMovieIdFromRequest(r)
		if err != nil {
			badRequest(w, "Invalid ID")
			return
		}

		_, err = db.DeleteVote(movieId)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		rowsAffected, err := db.DeleteMovie(movieId)
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		if rowsAffected == 0 {
			badRequest(w, "Not Found")
			return
		}

		success(w, nil)
	}
}

// POST /movies/{id}/vote
func VoteMovieHandler(db *database.DbModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authRequest(db, r)
		if err != nil {
			notAuthorized(w, err.Error())
			return
		}

		movieId, err := getMovieIdFromRequest(r)
		if err != nil {
			badRequest(w, "Invalid movie ID")
			return
		}

		params := r.URL.Query()
		if _, ok := params["score"]; !ok {
			badRequest(w, "Missing parameter: 'score'")
			return
		}

		score, err := strconv.ParseInt(params["score"][0], 10, 32)
		if err != nil {
			badRequest(w, "Invalid score value")
			return
		}

		if score < 0 || score > 10 {
			badRequest(w, "Score value outside acceptable range (0-10)")
			return
		}

		user, _ := getCurrentUser(db, r)
		err = db.VoteMovie(movieId, user.UserId, int(score))
		if err != nil {
			log.Printf("Database error: %s", err)
			internalError(w, "Problem with database")
			return
		}

		success(w, nil)
	}
}

// GET /goroutine-example
func LongOperationExampleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup

		response := ""
		worker := func(wnum int) {
			rand.Seed(time.Now().UnixNano())
			duration := 200 + rand.Intn(300)
			time.Sleep(time.Millisecond * time.Duration(duration))
			response = response + fmt.Sprintf("Goroutine #%d completed in %d ms\n", wnum, duration)
			wg.Done()
		}

		for i := 0; i < 3; i++ {
			wg.Add(1)
			go worker(i + 1)
		}

		wg.Wait()

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, response)
	}
}

func getCurrentUser(db *database.DbModel, r *http.Request) (*model.User, error) {
	token := r.Header.Get("token")
	if token == "" {
		return nil, errors.New("Missing token")
	}

	user, err := db.GetUserByToken(token)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("Invalid token")
	}

	return user, nil
}

func authRequest(db *database.DbModel, r *http.Request) error {
	_, err := getCurrentUser(db, r)
	return err
}

func getMovieIdFromRequest(r *http.Request) (int64, error) {
	params := mux.Vars(r)
	movieId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || movieId <= 0 {
		return 0, err
	}

	return movieId, nil
}

func validateMovieRequestData(movie *model.Movie) error {
	if movie.Title == "" {
		return errors.New("Field 'title' cannot be empty")
	}

	if len(movie.Title) > maxTitleLength {
		return errors.New("Field 'title' too long")
	}

	if len(movie.Description) > maxDescriptionLength {
		return errors.New("Field 'description' too long")
	}

	if movie.ImageURL == "" {
		return errors.New("Field 'image_url' cannot be empty")
	}

	if len(movie.ImageURL) > maxImageUrlLength {
		return errors.New("Field 'image_url' too long")
	}

	return nil
}

func getMovieDataFromRequest(r *http.Request) (*model.Movie, error) {
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		return nil, errors.New("Invalid JSON request body")
	}

	movie.Title = strings.TrimSpace(movie.Title)
	movie.Description = strings.TrimSpace(movie.Description)
	movie.ImageURL = strings.TrimSpace(movie.ImageURL)
	
	return &movie, nil
}
