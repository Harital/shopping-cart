package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver anonymously in order to make it work
)

func InitMySqlDB() (*sql.DB, error) {
	// All these values should be in a secret manarger storage. Hardcoded for simplicity
	user := "shopping-cart-app"
	password := "shoppingCartPassword!"
	host := "shopping-cart-mysql"
	port := 3306
	dataBase := "shoppingCart"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dataBase)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
