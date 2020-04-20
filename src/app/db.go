package main

import (
	"os"
	"fmt"
	"github.com/fatih/color"
	"database/sql"
	_ "github.com/lib/pq"
)

func dbWrite(product Clothing) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
    os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("USER"),
		os.Getenv("DBNAME"))

	db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }

	sqlStatement := `
	INSERT INTO floryday (product, code, description, price)
	VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, product.Name, product.Code, product.Description, product.Price)

	if err != nil {
		color.Red("[DB] Failed Write: %s", product.Name)
	  panic(err)
	}
	color.Green("[DB] Successful Write: %s", product.Name)
}
