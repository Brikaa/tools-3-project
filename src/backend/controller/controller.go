package controller

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Brikaa/tools-3-project/src/backend/repo"
	g "github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	db  *sql.DB
	rdb *redis.Client
}

func CreateController(db *sql.DB, rdb *redis.Client) Controller {
	return Controller{db, rdb}
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

type PutSlotRequest struct {
	Start time.Time `json:"start" time_format:"RFC3339"`
	End   time.Time `json:"end" time_format:"RFC3339"`
}

type PutAppointmentRequest struct {
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
	log.Print(*err)
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

var allowedRoles = map[string]bool{"doctor": true, "patient": true}
var isAlNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (controller Controller) Auth(role string, fn func(*UserContext, *g.Context)) func(*g.Context) {
	return func(ctx *g.Context) {
		authHeader := ctx.GetHeader("Authorization")
		var token string
		if authHeader != "" {
			authData := strings.Split(authHeader, " ")
			if len(authData) != 2 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("Invalid authorization header"))
				return
			}

			scheme := authData[0]
			if !strings.EqualFold(scheme, "Basic") {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("Invalid authorization scheme"))
				return
			}
			token = authData[1]
		} else {
			token = ctx.Query("token")
		}

		userpass, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			handleInternalServerError(ctx, &err)
			return
		}
		userpassData := strings.Split(string(userpass), ":")
		if len(userpassData) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("Invalid authorization token"))
			return
		}
		username := userpassData[0]
		password := userpassData[1]

		user, dbErr := repo.GetUserByIdAndPassword(controller.db, username, password)
		if dbErr != nil {
			handleInternalServerError(ctx, &dbErr)
			return
		}
		if user == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if role != "*" && user.Role != role {
			log.Print(user.Role)
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

func (controller Controller) withPutSlotBusinessRules(
	userCtx *UserContext, ctx *g.Context, fn func(*PutSlotRequest),
) {
	var req PutSlotRequest
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
			errorResponse(fmt.Sprintf("Slot with id %s overlaps with this configuration", *overlap)),
		)
		return
	}
	fn(&req)
}

func (controller Controller) CreateSlot(userCtx *UserContext, ctx *g.Context) {
	controller.withPutSlotBusinessRules(userCtx, ctx, func(req *PutSlotRequest) {
		if err := repo.InsertSlot(controller.db, req.Start, req.End, userCtx.ID); err != nil {
			handleInternalServerError(ctx, &err)
			return
		}
		ctx.Status(http.StatusCreated)
	})
}

func (controller Controller) UpdateSlot(userCtx *UserContext, ctx *g.Context) {
	controller.withPutSlotBusinessRules(userCtx, ctx, func(req *PutSlotRequest) {
		updated, err := repo.UpdateSlotByIdAndDoctorId(controller.db, ctx.Param("id"), userCtx.ID, req.Start, req.End)
		if err != nil {
			handleInternalServerError(ctx, &err)
			return
		}
		if !updated {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusOK)
	})
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

func (controller Controller) withPutAppointmentBusinessRules(
	userCtx *UserContext, ctx *g.Context, fn func(*PutAppointmentRequest),
) {
	var req PutAppointmentRequest
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
		handleInternalServerError(ctx, &targetErr)
		return
	}
	if target == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse("This slot does not exist"))
		return
	}

	fn(&req)
}

func createChannelName(doctorId string) string {
	return fmt.Sprintf("appointment.%s", doctorId)
}

type Message struct {
	DoctorId  string `json:"doctorId"`
	PatientId string `json:"patientId"`
	Operation string `json:"operation"`
}

func createMessage(doctorId, patientId, operation string) ([]byte, error) {
	return json.Marshal(Message{DoctorId: doctorId, PatientId: patientId, Operation: operation})
}

func (controller *Controller) publishMessage(
	ctx *g.Context, patientId, operation string, getDoctorId func() (*string, error),
) {
	doctorId, err := getDoctorId()
	if err != nil {
		log.Print(err)
		return
	}
	message, err := createMessage(*doctorId, patientId, operation)
	if err != nil {
		log.Print(err)
		return
	}
	if err := controller.rdb.Publish(ctx, createChannelName(*doctorId), message).Err(); err != nil {
		log.Print(err)
		return
	}
}

func (controller Controller) CreateAppointment(userCtx *UserContext, ctx *g.Context) {
	controller.withPutAppointmentBusinessRules(userCtx, ctx, func(req *PutAppointmentRequest) {
		if err := repo.InsertAppointment(controller.db, req.SlotID, userCtx.ID); err != nil {
			handleInternalServerError(ctx, &err)
			return
		}

		ctx.Status(http.StatusCreated)
		// Called ReservationCreated in the instructions
		controller.publishMessage(
			ctx,
			userCtx.ID,
			"AppointmentCreated",
			func() (*string, error) { return repo.GetDoctorIdBySlotId(controller.db, req.SlotID) },
		)
	})
}

func (controller Controller) UpdateAppointment(userCtx *UserContext, ctx *g.Context) {
	controller.withPutAppointmentBusinessRules(userCtx, ctx, func(req *PutAppointmentRequest) {
		prevDoctorId, err := repo.GetDoctorIdByAppointmentId(controller.db, ctx.Param("id"))
		if err != nil {
			handleInternalServerError(ctx, &err)
			return
		}
		updated, err := repo.UpdateAppointmentByIdAndPatientId(controller.db, ctx.Param("id"), userCtx.ID, req.SlotID)
		if err != nil {
			handleInternalServerError(ctx, &err)
			return
		}
		if !updated {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusOK)
		// Called ReservationUpdated in the instructions
		controller.publishMessage(
			ctx,
			userCtx.ID,
			"AppointmentUpdated",
			func() (*string, error) { return repo.GetDoctorIdBySlotId(controller.db, req.SlotID) },
		)
		controller.publishMessage(
			ctx,
			userCtx.ID,
			"AppointmentUpdated",
			func() (*string, error) { return prevDoctorId, nil },
		)
	})
}

func (controller Controller) DeleteAppointment(userCtx *UserContext, ctx *g.Context) {
	doctorId, doctorErr := repo.GetDoctorIdByAppointmentId(controller.db, ctx.Param("id"))
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
	// Called ReservationCancelled in the instructions
	controller.publishMessage(
		ctx,
		userCtx.ID,
		"AppointmentCancelled",
		func() (*string, error) { return doctorId, doctorErr },
	)
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

func (controller Controller) GetCurrentUser(userContext *UserContext, ctx *g.Context) {
	user, err := repo.GetUserById(controller.db, userContext.ID)
	if err != nil {
		handleInternalServerError(ctx, &err)
		return
	}
	if user == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.IndentedJSON(http.StatusOK, g.H{"user": user})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readMessages(ws *websocket.Conn, ch <-chan *redis.Message) {
	for message := range ch {
		ws.WriteMessage(websocket.TextMessage, []byte(message.Payload))
	}
}

func (controller *Controller) GetAppointmentUpdates(userCtx *UserContext, ctx *g.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	channelName := createChannelName(userCtx.ID)
	pubsub := controller.rdb.Subscribe(ctx, channelName)
	defer pubsub.Close()
	ch := pubsub.Channel()
	go readMessages(ws, ch)
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			pubsub.Unsubscribe(ctx, channelName)
			return
		}
	}
}
