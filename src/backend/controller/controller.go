package controller

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

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

type CreateSlotRequest struct {
	Start time.Time `json:"start" time_format:"RFC3339"`
	End   time.Time `json:"end" time_format:"RFC3339"`
}

type CreateAppointmentRequest struct {
	SlotID string `json:"slotId"`
}

type UserContext struct {
	ID   string
	Role string
}

func createUserContext(id string, role string) *UserContext {
	return &UserContext{ID: id, Role: role}
}

func errorResponse(message string) *g.H {
	return &g.H{"message": message}
}

func handleInternalServerError(ctx *g.Context, err *error) {
	log.Print(err)
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

var allowedRoles = map[string]bool{"doctor": true, "patient": true}
var isAlNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (controller Controller) Auth(role string, fn func(*UserContext, *g.Context)) func(*g.Context) {
	return func(ctx *g.Context) {
		authHeader := ctx.GetHeader("Authorization")
		authData := strings.Split(authHeader, " ")
		if len(authData) != 2 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("Invalid authorization header"))
			return
		}

		scheme := authData[0]
		token := authData[1]
		if !strings.EqualFold(scheme, "Basic") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("Invalid authorization scheme"))
			return
		}

		userpass, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			handleInternalServerError(ctx, &err)
			return
		}
		userpassData := strings.Split(string(userpass), ":")
		username := userpassData[0]
		password := userpassData[1]

		user, dbErr := repo.GetUserByUsernameAndPassword(controller.db, username, password)
		if dbErr != nil {
			handleInternalServerError(ctx, &dbErr)
			return
		}
		if user == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if role != "*" && user.Role != role {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		userContext := createUserContext(user.ID, user.Role)
		fn(userContext, ctx)
	}
}

func (controller Controller) Signup(ctx *g.Context) {
	var req SignUpRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	if _, ok := allowedRoles[req.Role]; !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("Invalid role"))
		return
	}
	if !isAlNum(req.Username) {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorResponse("Username can only contain alphabetic or numeric characters"),
		)
		return
	}
	if len(req.Username) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("Username must be non-empty"))
		return
	}
	user, err := repo.GetUserByUsername(controller.db, req.Username)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if user != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("A user with this username already exists"))
		return
	}
	if err := repo.InsertUser(controller.db, req.Username, req.Password, req.Role); err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (controller Controller) Login(ctx *g.Context) {
	var req LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	user, err := repo.GetUserByUsernameAndPassword(controller.db, req.Username, req.Password)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("Invalid username or password"))
		return
	}
	ctx.IndentedJSON(http.StatusOK,
		g.H{"token": base64.StdEncoding.EncodeToString([]byte(user.ID + ":" + user.Password))})
}

func (controller Controller) CreateSlot(userCtx *UserContext, ctx *g.Context) {
	var req CreateSlotRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	overlap, err := repo.GetOverlappingSlotId(controller.db, userCtx.ID, req.Start, req.End)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if !req.Start.Before(req.End) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("The slot must start before it ends"))
		return
	}
	if req.Start.Before(time.Now()) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("The slot must start in the future"))
		return
	}
	if overlap != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorResponse(fmt.Sprintf("Slot with id %v overlaps with this configuration", *overlap)),
		)
		return
	}
	if err := repo.InsertSlot(controller.db, req.Start, req.End, userCtx.ID); err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (controller Controller) DeleteSlot(userCtx *UserContext, ctx *g.Context) {
	deleted, err := repo.DeleteSlotByIdAndDoctorId(controller.db, ctx.Param("id"), userCtx.ID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if !deleted {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.Status(http.StatusOK)
}

func (controller Controller) GetSlots(userCtx *UserContext, ctx *g.Context) {
	slots, err := repo.GetSlotsByDoctorId(controller.db, userCtx.ID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, g.H{"slots": slots})
}

func (controller Controller) GetDoctorAppointments(userCtx *UserContext, ctx *g.Context) {
	appointments, err := repo.GetAppointmentsByDoctorId(controller.db, userCtx.ID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, g.H{"appointments": appointments})
}

func (controller Controller) GetAppointments(userCtx *UserContext, ctx *g.Context) {
	appointments, err := repo.GetAppointmentsByPatientId(controller.db, userCtx.ID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, g.H{"appointments": appointments})
}

func (controller Controller) CreateAppointment(userCtx *UserContext, ctx *g.Context) {
	var req CreateAppointmentRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	reserved, err := repo.GetAppointmentIdBySlotId(controller.db, req.SlotID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if reserved != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("This slot is already reserved"))
		return
	}

	target, targetErr := repo.GetSlotIdBySlotId(controller.db, req.SlotID)
	if targetErr != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if target == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("This slot does not exist"))
		return
	}

	if err := repo.InsertAppointment(controller.db, req.SlotID, userCtx.ID); err != nil {
		handleInternalServerError(ctx, &err)
	}

	ctx.Status(http.StatusOK)
}

func (controller Controller) DeleteAppointment(userCtx *UserContext, ctx *g.Context) {
	deleted, err := repo.DeleteAppointmentByIdAndPatientId(controller.db, ctx.Param("id"), userCtx.ID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if !deleted {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.Status(http.StatusOK)
}

func (controller Controller) GetDoctors(_ *UserContext, ctx *g.Context) {
	doctors, err := repo.GetDoctors(controller.db)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, g.H{"doctors": doctors})
}

func (controller Controller) GetAvailableSlotsForDoctor(_ *UserContext, ctx *g.Context) {
	slots, err := repo.GetAvailableSlotsByDoctorId(controller.db, ctx.Param("id"))
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, g.H{"slots": slots})
}
