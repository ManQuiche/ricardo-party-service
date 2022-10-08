package party

import (
	"errors"
	"github.com/gin-gonic/gin"
	ricardoErr "gitlab.com/ricardo-public/errors/pkg/errors"
	tokens "gitlab.com/ricardo-public/jwt-tools/v2/pkg/token"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
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
	var cpr entities.CreatePartyRequest
	err := gtx.ShouldBindJSON(&cpr)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, err.Error()))
		return
	}
	userID, _ := gtx.Get(tokens.UserIDKey)

	p := entities.Party{
		Name:        cpr.Name,
		UserID:      userID.(uint),
		Time:        cpr.Time,
		Description: cpr.Description,
	}
	cParty, err := c.service.Save(gtx.Request.Context(), p)
	if err != nil {
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
// @Router /parties/{party_id} [PATCH]
func (c controller) Update(gtx *gin.Context) {
	var upr entities.UpdatePartyRequest
	err := gtx.ShouldBindJSON(&upr)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, err.Error()))
		return
	}

	partyId, err := strconv.Atoi(gtx.Param("party_id"))
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, "invalid ID format"))
		return
	}
	uintPartyId := uint(partyId)

	_, err = c.canUpdateOrDelete(gtx, uintPartyId)
	if err != nil {
		return
	}

	p := entities.Party{
		ID:     uintPartyId,
		Name:   upr.Name,
		UserID: upr.UserID,
	}

	uParty, err := c.service.Save(gtx.Request.Context(), p)
	if err != nil {
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
	log.Println(gtx.Get("sub"))

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
	partyId, err := strconv.ParseUint(gtx.Param("party_id"), 10, 32)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, "invalid ID format"))
		return
	}
	uintPartyId := uint(partyId)

	_, err = c.canUpdateOrDelete(gtx, uintPartyId)
	if err != nil {
		return
	}

	err = c.service.Delete(gtx.Request.Context(), uintPartyId)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return
	}

	gtx.Status(http.StatusOK)
}

func (c controller) canUpdateOrDelete(gtx *gin.Context, partyID uint) (bool, error) {
	p, err := c.service.Get(gtx.Request.Context(), partyID)
	if err != nil {
		_ = ricardoErr.GinErrorHandler(gtx, ricardoErr.New(ricardoErr.ErrBadRequest, ""))
		return false, err
	}

	userID, _ := gtx.Get(tokens.UserIDKey)
	pToken, _ := tokens.ParseFromGinContext(gtx, c.accessSecret)
	claims, _ := tokens.Claims(pToken)

	if userID.(uint) != p.UserID && tokens.IsAdmin(*claims) {
		err = errors.New("unauthorized to update or delete")
		_ = ricardoErr.GinErrorHandler(gtx, err)
		return false, err
	}

	return true, nil
}
