package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Brikaa/tools-3-project/src/backend/model"
	"github.com/Brikaa/tools-3-project/src/backend/repo"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	db *sql.DB
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func errorResponse(message string) gin.H {
	return gin.H{"message": message}
}

var allowedRoles = map[string]bool{"doctor": true, "patient": true}

func CreateController(db *sql.DB) Controller {
	return Controller{db}
}
func (controller Controller) Signup(c *gin.Context) {
	var req SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}
	if _, ok := allowedRoles[req.Role]; !ok {
		log.Print(req)
		c.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid role"))
		return
	}
	user, err := repo.SelectUserByUsername(controller.db, req.Username)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if user != nil {
		c.IndentedJSON(http.StatusBadRequest, errorResponse("A user with this username already exists"))
	}
	newUser := model.User{Username: req.Username, Password: req.Password, Role: req.Role}
	if err := repo.InsertUser(controller.db, newUser); err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}
