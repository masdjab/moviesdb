package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func RunMigration(database *sql.DB) {
	log.Println("Starting database migration...")
	filenames := getMigrationFileNames("database/migrations")
	for _, fn := range filenames {
		log.Println(fn)
		query, err := ioutil.ReadFile(fn)
		if err != nil {
			panic(err)
		}

		splits := strings.Split(string(query), ";\n")
		for _, query := range splits {
			query = strings.TrimSpace(query)
			if query != "" {
				query = fmt.Sprintf("%s;", query)
				_, err = database.Exec(query)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	log.Println("Database migration finished")
}

func getMigrationFileNames(dir string) []string {
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	filenames := []string{}
	for _, fi := range fileInfo {
		fn := fi.Name()
		offsetToExt := len(fn) - 4
		if fn[offsetToExt:] == ".sql" {
			filenames = append(filenames, fmt.Sprintf("%s/%s", dir, fn))
		}
	}

	sort.Strings(filenames)

	return filenames
}
