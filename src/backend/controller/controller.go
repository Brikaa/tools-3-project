package controller

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Brikaa/tools-3-project/src/backend/model"
	"github.com/Brikaa/tools-3-project/src/backend/repo"
	g "github.com/gin-gonic/gin"
)

type Controller struct {
	db *sql.DB
}

func CreateController(db *sql.DB) Controller {
	return Controller{db}
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

type UserContext struct {
	Id   string
	Role string
}

func createUserContext(id string, role string) *UserContext {
	return &UserContext{Id: id, Role: role}
}

func errorResponse(message string) *g.H {
	return &g.H{"message": message}
}

var allowedRoles = map[string]bool{"doctor": true, "patient": true}
var isAlNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (controller Controller) auth(role string) func(*g.Context) {
	return func(c *g.Context) {
		authHeader := c.GetHeader("Authorization")
		authData := strings.Split(authHeader, " ")
		if len(authData) != 2 {
			c.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid authorization header"))
			return
		}

		scheme := authData[0]
		token := authData[1]
		if !strings.EqualFold(scheme, "Basic") {
			c.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid authorization scheme"))
			return
		}

		userpass, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		userpassData := strings.Split(string(userpass), ":")
		username := userpassData[0]
		password := userpassData[1]

		user, dbErr := repo.SelectUserByUsernameAndPassword(controller.db, username, password)
		if dbErr != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if user == nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		if role != "*" && user.Role != role {
			c.Status(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func (controller Controller) Signup(ctx *g.Context) {
	var req SignUpRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	if _, ok := allowedRoles[req.Role]; !ok {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid role"))
		return
	}
	if !isAlNum(req.Username) {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse("Username can only contain alphabetic or numeric characters"))
		return
	}
	if len(req.Username) == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse("Username must be non-empty"))
		return
	}
	user, err := repo.SelectUserByUsername(controller.db, req.Username)
	if err != nil {
		log.Print(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if user != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse("A user with this username already exists"))
	}
	newUser := model.CreateUser(req.Username, req.Password, req.Role)
	if err := repo.InsertUser(controller.db, newUser); err != nil {
		log.Print(err)
		ctx.Status(http.StatusInternalServerError)
	}
	ctx.Status(http.StatusOK)
}

func (controller Controller) Login(ctx *g.Context) {
	var req LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	user, err := repo.SelectUserByUsernameAndPassword(controller.db, req.Username, req.Password)
	if err != nil {
		log.Print(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse("Invalid username or password"))
		return
	}
	ctx.IndentedJSON(http.StatusOK,
		g.H{"token": base64.StdEncoding.EncodeToString([]byte(user.ID + ":" + user.Password))})
	return
}
