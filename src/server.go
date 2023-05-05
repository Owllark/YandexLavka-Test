package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"yandex-team.ru/bstask/db"
	"yandex-team.ru/bstask/routes"
)

func main() {
	//e := setupServer()
	//e.Logger.Fatal(e.Start(":8080"))
	var db = new(db.PostgreSQLDatabase)
	err := db.Connect()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected Successfully!")
	}
	_, err = db.Exec("SELECT * FROM test")
	if err != nil {
		fmt.Println(err)
	} else {
		rows, _ := db.Query("SELECT idx FROM test")
		for rows.Next() {
			var idx int
			if err := rows.Scan(&idx); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(idx)
			}

		}
	}

}

func setupServer() *echo.Echo {
	e := echo.New()
	routes.SetupRoutes(e)
	return e
}
