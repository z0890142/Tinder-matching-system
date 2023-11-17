package restful

import (
	"errors"
	"net/http"
	"strconv"
	"tinderMatchingSystem/internal/models"

	"github.com/gin-gonic/gin"

	httpErr "tinderMatchingSystem/internal/error"
	matcheSystem "tinderMatchingSystem/pkg/domain/match_system"
)

// @title Tinder Matching System API
// @version 1.0
type restfulHandler struct {
	matcheSystem matcheSystem.MatcheSystem
}

func NewRestfulHandler() restfulHandler {
	matchesystem := matcheSystem.NewMatcheSystem()
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

	matches = h.matcheSystem.MatchingPerson(&person, person.NumberOfDate)

	response := models.AddSinglePersonAndMatchResponse{
		NewUser: &person,
		Matches: matches,
	}

	ginC.JSON(http.StatusOK, response)

}

// @Summary Remove a user from the matching system so that the user cannot be matched anymore.
// @router /tinder_system/v1/persons/{name} [delete]
// @param name path string true "person name"
// @Success 200 {object} models.CommonResponse
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
	ginC.JSON(http.StatusOK, models.CommonResponse{
		Result: "OK",
	})
}

// @Summary Find the most N possible matched single people, where N is a request parameter.
// @router /tinder_system/v1/persons [get]
// @param n query int true "number of single person"
// @param name query string false "person name"
// @param gender query string false "M or F"
// @param heightGte query int false "person height >="
// @param heightLte query string false "person height <="
// @Success 200 {object} models.QuerySinglePeopleResponse
// @Failure 400 {object} models.ErrorResponse
func (h *restfulHandler) QuerySinglePeople(ginC *gin.Context) {

	queryFilter, err := getFilter(ginC)
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}

	matches, err := h.matcheSystem.QuerySinglePerson(queryFilter)
	if err != nil {
		err := httpErr.NewHttpError(http.StatusBadRequest, httpErr.MatcheSinglePersonErr, err.Error())
		ginC.Error(err)
		return
	}
	ginC.JSON(http.StatusOK, models.QuerySinglePeopleResponse{
		Matches: matches,
	})
}

func getFilter(ginC *gin.Context) (*models.QueryFilter, error) {
	heightGte, err := parseInt(ginC, "heightGte")
	if err != nil {
		return nil, err
	}

	heightLte, err := parseInt(ginC, "heightLte")
	if err != nil {
		return nil, err
	}

	n, err := parseInt(ginC, "n")
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, err
	}

	return &models.QueryFilter{
		Name:      ginC.Query("name"),
		Gender:    ginC.Query("gender"),
		HeightGte: heightGte,
		HeightLte: heightLte,
		N:         n,
	}, nil
}

func parseInt(ginC *gin.Context, paramName string) (int, error) {
	valStr := ginC.Query(paramName)
	if valStr == "" {
		return 0, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, errors.New(paramName + " must be a valid integer")
	}
	return val, nil
}
