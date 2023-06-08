package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitMysql(url string) *sql.DB {
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("mysql is  connected at %s\n", url)

	return db
}
