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
	cfg := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
		DBName:    "app",
		ParseTime: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully pinged the DBMS")

	router := g.Default()
	router.GET("/health", health)

	controller := controller.CreateController(db)
	router.POST("/signup", controller.Signup)
	router.POST("/login", controller.Login)

	router.GET("/user", controller.Auth("*", controller.GetCurrentUser))

	router.PUT("/slots", controller.Auth("doctor", controller.CreateSlot))
	router.PUT("/slots/:id", controller.Auth("doctor", controller.UpdateSlot))
	router.DELETE("/slots/:id", controller.Auth("doctor", controller.DeleteSlot))
	router.GET("/slots", controller.Auth("doctor", controller.GetSlots))
	router.GET("/doctor-appointments", controller.Auth("doctor", controller.GetDoctorAppointments))

	router.GET("/appointments", controller.Auth("patient", controller.GetAppointments))
	router.PUT("/appointments", controller.Auth("patient", controller.CreateAppointment))
	router.PUT("/appointments/:id", controller.Auth("patient", controller.UpdateAppointment))
	router.DELETE("/appointments/:id", controller.Auth("patient", controller.DeleteAppointment))
	router.GET("/doctors", controller.Auth("patient", controller.GetDoctors))
	router.GET("/doctors/:id/slots", controller.Auth("patient", controller.GetAvailableSlotsForDoctor))

	router.Run(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")))
}
