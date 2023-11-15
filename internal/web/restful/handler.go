package restful

import (
	"net/http"
	"strconv"
	httpErr "tinderMatchingSystem/error"
	"tinderMatchingSystem/internal/models"
	matchesystem "tinderMatchingSystem/pkg/domain/match_system"

	"github.com/gin-gonic/gin"
)

// @title Tinder Matching System API
// @version 1.0
type restfulHandler struct {
	matcheSystem matchesystem.MatcheSystem
}

func NewRestfulHandler() restfulHandler {
	matchesystem := matchesystem.NewMatcheSystem()
	return restfulHandler{
		matcheSystem: matchesystem,
	}
}

// @Summary Add a new user to the matching system and find any possible matches for the new user
// @router /tinder_system/v1/persons [post]
// @param params body models.SinglePerson true "params"
// @Success 200 {object} models.AddSinglePersonAndMatchResponse
// @Failure 400 {object} models.ErrorResponse
func (h *restfulHandler) AddSinglePersonAndMatch(ginC *gin.Context) {
	person := models.SinglePerson{}
	err := ginC.BindJSON(&person)
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.BindJSONErr, err.Error())
		ginC.Error(err)
		return
	}

	person.Lock.Lock()
	defer person.Lock.Unlock()

	err = h.matcheSystem.RegisterSinglePerson(&person)
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.RemoveSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}

	var matches []*models.SinglePerson

	if person.Gender == "M" {
		matches = h.matcheSystem.MatchingFmalePerson(&person, person.NumberOfDate)
	} else if person.Gender == "F" {
		matches = h.matcheSystem.MatchingMalePerson(&person, person.NumberOfDate)
	}

	response := models.AddSinglePersonAndMatchResponse{
		NewUser: &person,
		Matches: matches,
	}

	ginC.JSON(http.StatusOK, response)

}

// @Summary Remove a user from the matching system so that the user cannot be matched anymore.
// @router /tinder_system/v1/persons/{name} [delete]
// @param name path string true "person name"
// @Success 200 {string} ok
// @Failure 400 {object} models.ErrorResponse
func (h *restfulHandler) RemoveSinglePerson(ginC *gin.Context) {
	name, exist := ginC.Params.Get("name")
	if !exist {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.GetParamsErr, "params miss")
		ginC.Error(err)
		return
	}

	err := h.matcheSystem.RemoveSinglePerson(name)
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.RemoveSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}
	ginC.String(http.StatusOK, "OK")
}

// @Summary Find the most N possible matched single people, where N is a request parameter.
// @router /tinder_system/v1/persons [get]
// @param n query int true "number of single person"
// @param name query string false "person name"
// @param gender query string false "M or F"
// @param heightGte query int false "person height >="
// @param heightLte query string false "person height <="
// @Success 200 {string} ok
// @Failure 400 {object} models.ErrorResponse
func (h *restfulHandler) QuerySinglePeople(ginC *gin.Context) {
	heightGte, err := strconv.Atoi(ginC.Query("heightGte"))
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}
	heightLte, err := strconv.Atoi(ginC.Query("heightLte"))
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}
	n, err := strconv.Atoi(ginC.Query("n"))
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}
	if n == 0 {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, "n can't null")
		ginC.Error(err)
		return
	}

	queryFilter := models.QueryFilter{
		Name:      ginC.Query("name"),
		Gender:    ginC.Query("gender"),
		HeightGte: heightGte,
		HeightLte: heightLte,
		N:         n,
	}

	matches, err := h.matcheSystem.QuerySinglePerson(queryFilter)
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}
	ginC.JSON(http.StatusOK, matches)
}
