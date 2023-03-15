package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"moviesdb.com/config"
)

type DbModel struct {
	conn *sql.DB
}


func Connect(conf config.DatabaseConfig) (*sql.DB, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Name)
	conn, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewDbModel(conn *sql.DB) *DbModel {
	return &DbModel{conn: conn}
}

func escape(cmd string) string {
	dest := make([]byte, 0, 2*len(cmd))
	var escape byte
	for i := 0; i < len(cmd); i++ {
		c := cmd[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}
