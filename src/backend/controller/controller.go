package controller

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"regexp"

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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func errorResponse(message string) gin.H {
	return gin.H{"message": message}
}

var allowedRoles = map[string]bool{"doctor": true, "patient": true}
var isAlNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func CreateController(db *sql.DB) Controller {
	return Controller{db}
}

func (controller Controller) Signup(c *gin.Context) {
	var req SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}
	if _, ok := allowedRoles[req.Role]; !ok {
		c.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid role"))
		return
	}
	if !isAlNum(req.Username) {
		c.IndentedJSON(http.StatusBadRequest, errorResponse("Username can only contain alphabetic or numeric characters"))
		return
	}
	if len(req.Username) == 0 {
		c.IndentedJSON(http.StatusBadRequest, errorResponse("Username must be non-empty"))
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

func (controller Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}
	user, err := repo.SelectUserByUsernameAndPassword(controller.db, req.Username, req.Password)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if user == nil {
		c.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid username or password"))
		return
	}
	c.IndentedJSON(http.StatusOK,
		gin.H{"token": base64.StdEncoding.EncodeToString([]byte(req.Username + ":" + req.Password))})
	return
}
