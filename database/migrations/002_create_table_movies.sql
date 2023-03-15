DROP TABLE IF EXISTS movies;

CREATE TABLE movies(
  id            INT(1) UNSIGNED NOT NULL AUTO_INCREMENT, 
  title         VARCHAR(50) NOT NULL, 
  description   VARCHAR(200), 
  image_url     VARCHAR(200), 
  created_at    INT(1) UNSIGNED NOT NULL, 
  updated_at    INT(1) UNSIGNED NOT NULL,
  PRIMARY KEY(id)
);
