package bootstrap

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func InitDBConnection(user, password, host, port, database string) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	db, err := sql.Open("mysql",
		dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}
