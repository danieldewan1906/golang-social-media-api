package connection

import (
	"database/sql"
	"fmt"
	"golang-restful-api/internal/config"
	"log"

	_ "github.com/lib/pq"
)

func GetDatabase(conf config.Database) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable Timezone=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Pass,
		conf.Name,
		conf.Schema,
		conf.Tz,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open connection: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping connection: ", err.Error())
	}

	return db
}
