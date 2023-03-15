DROP TABLE IF EXISTS votes;

CREATE TABLE votes(
  movie_id      INT(1) UNSIGNED NOT NULL, 
  user_id       INT(1) UNSIGNED NOT NULL, 
  score         INT(1) NOT NULL, 
  created_at    INT(1) UNSIGNED NOT NULL, 
  updated_at    INT(1) UNSIGNED NOT NULL,
  KEY movie_id(movie_id)
);
