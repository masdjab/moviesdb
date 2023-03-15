package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"moviesdb.com/config"
	"moviesdb.com/database"
)

type Server struct {
	port int
	conn *sql.DB
}

func NewServer(conn *sql.DB) *Server {
	return &Server{
		port: config.ServerPort(),
		conn: conn,
	}
}

func (s *Server) createHandler() *mux.Router {
	db := database.NewDbModel(s.conn)
	router := mux.NewRouter()

	router.HandleFunc("/ping", PingHandler()).Methods("GET")
	router.HandleFunc("/login", LoginHandler(db)).Methods("POST")
	router.HandleFunc("/logout", LogoutHandler(db)).Methods("GET")
	router.HandleFunc("/movies", GetMoviesHandler(db)).Methods("GET")
	router.HandleFunc("/movies/{id}", GetMovieDetailHandler(db)).Methods("GET")
	router.HandleFunc("/movies", PostMovieHandler(db)).Methods("POST")
	router.HandleFunc("/movies/{id}", PatchMovieHandler(db)).Methods("PATCH")
	router.HandleFunc("/movies/{id}", DeleteMovieHandler(db)).Methods("DELETE")
	router.HandleFunc("/movies/{id}/vote", VoteMovieHandler(db)).Methods("POST")

	return router
}

func (s *Server) Start() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.createHandler(),
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	log.Printf("MoviesDB serving at port %d", s.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
