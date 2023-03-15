TRUNCATE users;

INSERT INTO users(uname, display_name, email, psswd, last_login) VALUES
('admin', 'Administrator', 'admin@moviesdb.com', 'e93e7bbef19e663ebb67c799f6a252b7', 0), 
('someuser', 'Some User', 'someuser@moviesdb.com', 'e93e7bbef19e663ebb67c799f6a252b7', 0), 
('anotheruser', 'Another User', 'anotheruser@moviesdb.com', 'e93e7bbef19e663ebb67c799f6a252b7', 0);
