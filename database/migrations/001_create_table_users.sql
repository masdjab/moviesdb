DROP TABLE IF EXISTS users;

CREATE TABLE users(
  id            INT(1) UNSIGNED NOT NULL AUTO_INCREMENT, 
  display_name  VARCHAR(50) NOT NULL, 
  email         VARCHAR(50) NOT NULL, 
  uname         VARCHAR(50) NOT NULL, 
  psswd         VARCHAR(32) NOT NULL, 
  token         VARCHAR(128) NOT NULL DEFAULT '', 
  last_login    INT(1) UNSIGNED NOT NULL, 
  PRIMARY KEY(id), 
  KEY email(email), 
  KEY uname(uname)
);
