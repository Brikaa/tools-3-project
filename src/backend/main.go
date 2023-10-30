package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Brikaa/tools-3-project/src/backend/controller"
	g "github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func health(c *g.Context) {
	c.String(http.StatusOK, "Up and running")
}

func main() {
	router := g.Default()
	cfg := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Net:       "tcp",
		Addr:      "database:3306",
		DBName:    "app",
		ParseTime: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	controller := controller.CreateController(db)
	fmt.Println("Successfully pinged the DBMS")
	router.GET("/health", health)
	router.POST("/signup", controller.Signup)
	router.POST("/login", controller.Login)
	router.PUT("/slots", controller.Auth("doctor", controller.InsertSlot))
	router.DELETE("/slots/:id", controller.Auth("doctor", controller.DeleteSlot))

	router.Run(":8000")
}
