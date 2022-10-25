package party

import (
	"errors"
	"github.com/gin-gonic/gin"
	ricardoErr "gitlab.com/ricardo-public/errors/pkg/errors"
	tokens "gitlab.com/ricardo-public/jwt-tools/v2/pkg/token"
	"gitlab.com/ricardo-public/tracing/pkg/tracing"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strconv"
)

type Controller interface {
	Create(gtx *gin.Context)
	Update(gtx *gin.Context)
	Get(gtx *gin.Context)
	GetForUser(gtx *gin.Context)
	GetOne(gtx *gin.Context)
	Delete(gtx *gin.Context)
}

type controller struct {
	service      app.PartyService
	accessSecret []byte
}

func NewController(service app.PartyService, accessSecret []byte) Controller {
	return controller{service: service, accessSecret: accessSecret}
}

// Create
// @Summary Create a party
// @Description Create a party from json
// @Param party body entities.CreatePartyRequest true "Add party"
// @Success 200 {object} entities.Party
// @Failure 400 {object} ricardoErr.RicardoError
// @Router /parties [post]
func (c controller) Create(gtx *gin.Context) {
	span := gtx.Request.Context().Value(tracing.HttpSpanKey).(trace.Span)
	var err error
	defer span.End()

	var cpr entities.CreatePartyRequest
	err = gtx.ShouldBindJSON(&cpr)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusBadRequest)))
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, err.Error()))
		return
	}

	userID, ok := gtx.Get(tokens.UserIDKey)
	if !ok {
		err = errors.New("cannot retrieve token from gin context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusBadRequest)))
		log.Println(err)
	}

	p := entities.Party{
		Name:        cpr.Name,
		UserID:      userID.(uint),
		Time:        cpr.Time,
		Description: cpr.Description,
		Location:    cpr.Location,
	}
	cParty, err := c.service.Save(gtx.Request.Context(), p)
	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusInternalServerError)))
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, cParty)
}

// Update
// @Summary Update a party
// @Description Update a party from json
// @Param party_id path int true "Party id"
// @Param party body entities.UpdatePartyRequest true "Updated party info"
// @Success 200 {object} entities.Party
// @Failure 400 {object} ricardoErr.RicardoError
// @Router /parties/{party_id} [PUT]
func (c controller) Update(gtx *gin.Context) {
	span := gtx.Request.Context().Value(tracing.HttpSpanKey).(trace.Span)
	defer span.End()

	var upr entities.UpdatePartyRequest
	err := gtx.ShouldBindJSON(&upr)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusBadRequest)))
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, err.Error()))
		return
	}

	partyId, err := strconv.Atoi(gtx.Param("party_id"))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusBadRequest)))
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, "invalid ID format"))
		return
	}
	uintPartyId := uint(partyId)

	_, err = c.canUpdateOrDelete(gtx, uintPartyId, span)
	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusForbidden)))
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	parties, err := c.service.Get(gtx.Request.Context(), uint(partyId))
	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusNotFound)))
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	parties.Name = upr.Name
	parties.Description = upr.Description
	parties.Location = upr.Location

	uParty, err := c.service.Save(gtx.Request.Context(), *parties)
	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusInternalServerError)))
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, uParty)
}

// Get
// @Summary Get all party
// @Description Get all party
// @Success 200 {array} []entities.Party
// @Failure 400 {object} ricardoErr.RicardoError
// @Router /parties [GET]
func (c controller) Get(gtx *gin.Context) {
	parties, err := c.service.GetAll(gtx.Request.Context())
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, ""))
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

// GetForUser
// @Summary Get all parties of a user
// @Description Get all parties of a user
// @Param user_id path int true "User id"
// @Success 200 {array} []entities.Party
// @Failure 400 {object} ricardoErr.RicardoError
// @Router /parties/user/{user_id} [GET]
func (c controller) GetForUser(gtx *gin.Context) {
	userId, err := strconv.ParseUint(gtx.Param("user_id"), 10, 64)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, err.Error()))
	}

	parties, err := c.service.GetAllForUser(gtx.Request.Context(), uint(userId))
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

// GetOne
// @Summary Get a specific party by ID
// @Description Get a specific party by ID
// @Param party_id path int true "Party id"
// @Success 200 {object} entities.Party
// @Failure 400 {object} ricardoErr.RicardoError
// @Router /parties/{party_id} [GET]
func (c controller) GetOne(gtx *gin.Context) {
	partyId, err := strconv.ParseUint(gtx.Param("party_id"), 10, 64)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, ""))
		return
	}

	parties, err := c.service.Get(gtx.Request.Context(), uint(partyId))
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

// Delete
// @Summary Delete a party
// @Description Delete a party by id
// @Param party_id path int true "Party id"
// @Success 200
// @Failure 400 {object} ricardoErr.RicardoError
// @Router /parties/{party_id} [DELETE]
func (c controller) Delete(gtx *gin.Context) {
	span := gtx.Request.Context().Value(tracing.HttpSpanKey).(trace.Span)
	defer span.End()

	partyId, err := strconv.ParseUint(gtx.Param("party_id"), 10, 64)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusBadRequest)))
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, "invalid ID format"))
		return
	}
	uintPartyId := uint(partyId)

	_, err = c.canUpdateOrDelete(gtx, uintPartyId, span)
	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusForbidden)))
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	err = c.service.Delete(gtx.Request.Context(), uintPartyId)
	if err != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(http.StatusInternalServerError)))
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	gtx.Status(http.StatusOK)
}

func (c controller) canUpdateOrDelete(gtx *gin.Context, partyID uint, span trace.Span) (bool, error) {
	nctx, span := tracing.Tracer.Start(gtx.Request.Context(), "party.controller.canUpdateOrDelete")
	defer span.End()

	span.SetAttributes(attribute.Int("party.id", int(partyID)))

	p, err := c.service.Get(nctx, partyID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrNotFound, ""))
		return false, err
	}

	userID, _ := gtx.Get(tokens.UserIDKey)
	pToken, _ := tokens.ParseFromGinContext(gtx, c.accessSecret)
	claims, _ := tokens.Claims(pToken)

	span.SetAttributes(attribute.Int("user.id", int(userID.(uint))))

	if userID.(uint) != p.UserID && tokens.IsAdmin(*claims) {
		err = errors.New("unauthorized to update or delete")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return false, err
	}

	return true, nil
}
