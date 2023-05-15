package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"yandex-team.ru/bstask/routes"
)

func main() {
	e, err := setupServer()
	if err != nil {
		log.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":8080"))

}

const (
	DbUser     = "postgres"
	DbPassword = "password"
	DbName     = "lavka"
	DbHost     = "db"
)

func setupServer() (*echo.Echo, error) {
	e := echo.New()
	err := routes.ConnectWithDataBase(DbUser, DbPassword, DbName, DbHost)
	routes.SetupRoutes(e)
	return e, err
}
