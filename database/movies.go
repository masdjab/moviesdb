package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"moviesdb.com/model"
)

func generateSelectQuery(id int64, searchKeywords string) string {
	filters := []string{}

	if id > 0 {
		filters = append(filters, fmt.Sprintf("(id = %d)", id))
	}

	keywords := strings.Split(searchKeywords, " ")
	for _, k := range keywords {
		keyword := strings.Trim(k, " ")
		if len(keyword) > 0 {
			filters = append(filters, fmt.Sprintf("(title LIKE '%%%s%%')", keyword))
		}
	}

	cmd := `SELECT id, title, description, IFNULL(rating, 0), image_url, created_at, updated_at 
					FROM movies 
						LEFT JOIN (SELECT movie_id, AVG(score) AS rating FROM votes GROUP BY movie_id) 
						AS scores ON scores.movie_id = movies.id
					%s`
	
	whereCmd := strings.Join(filters, " AND ")
	if whereCmd != "" {
		whereCmd = " WHERE " + whereCmd
	}

	return fmt.Sprintf(cmd, whereCmd)
}

func (db *DbModel) GetMovies(keywords string) ([]model.Movie, error) {
	query := generateSelectQuery(0, keywords)
	results, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}

	var movies []model.Movie
	for results.Next() {
		var movie model.Movie
		var created int64
		var updated int64

		err = results.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Rating,
			&movie.ImageURL, &created, &updated)
		if err != nil {
			return nil, err
		}

		movie.CreatedAt = time.Unix(created, 0)
		movie.UpdatedAt = time.Unix(updated, 0)
		movies = append(movies, movie)
	}

	return movies, nil
}

func (db *DbModel) GetMovieDetail(id int64) (*model.Movie, error) {
	var movie model.Movie
	var created int64
	var updated int64

	cmd := generateSelectQuery(id, "")
	err := db.conn.QueryRow(cmd).Scan(&movie.ID, &movie.Title, &movie.Description,
		&movie.Rating, &movie.ImageURL, &created, &updated)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	movie.CreatedAt = time.Unix(created, 0)
	movie.UpdatedAt = time.Unix(updated, 0)

	return &movie, nil
}

func (db *DbModel) InsertMovie(movie *model.Movie) (int64, error) {
	cmd := `INSERT INTO movies(title, description, image_url, 
	        created_at, updated_at) VALUES('%s', '%s', '%s', 
	        UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`
	query := fmt.Sprintf(cmd,
		escape(movie.Title),
		escape(movie.Description),
		escape(movie.ImageURL))
	result, err := db.conn.Exec(query)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (db *DbModel) UpdateMovie(id int64, movie *model.Movie) error {
	cmd := `UPDATE movies SET title = '%s', description = '%s', 
	        image_url = '%s', updated_at = UNIX_TIMESTAMP() WHERE id = %d`
	query := fmt.Sprintf(cmd,
		escape(movie.Title),
		escape(movie.Description),
		escape(movie.ImageURL),
		id)
	_, err := db.conn.Query(query)
	return err
}

func (db *DbModel) DeleteMovie(id int64) (int64, error) {
	cmd := fmt.Sprintf("DELETE FROM movies WHERE id = %d", id)
	result, err := db.conn.Exec(cmd)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (db *DbModel) VoteMovie(movieId int64, userId int64, score int) error {
  voted, err := db.isMovieVoted(movieId, userId)
	if err != nil {
		return err
	}

	var query string
	if !voted {
		cmd := `INSERT INTO votes(movie_id, user_id, score, created_at, updated_at) 
			VALUES(%d, %d, %d, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`
		query = fmt.Sprintf(cmd, movieId, userId, score)
	} else {
		cmd := `UPDATE votes 
			SET score = %d, updated_at = UNIX_TIMESTAMP() 
			WHERE (movie_id = %d) AND (user_id = %d)`
		query = fmt.Sprintf(cmd, score, movieId, userId)
	}

	_, err = db.conn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *DbModel) DeleteVote(movieId int64) (int64, error) {
	cmd := fmt.Sprintf("DELETE FROM votes WHERE movie_id = %d", movieId)
	result, err := db.conn.Exec(cmd)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (db *DbModel) isMovieVoted(movieId int64, userId int64) (bool, error) {
	cmd := "SELECT movie_id FROM votes WHERE (movie_id = %d) AND (user_id = %d)"
	query := fmt.Sprintf(cmd, movieId, userId)
	err := db.conn.QueryRow(query).Scan(&movieId)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
