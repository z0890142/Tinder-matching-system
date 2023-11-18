package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"tinderMatchingSystem/internal/c"
	"tinderMatchingSystem/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite
	Listener    net.Listener
	HTTPServer  *httptest.Server
	testData    []*models.SinglePerson
	sharedState int
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func (suite *APISuite) SetupTest() {
	suite.setupGin()
	suite.setupData()

}

func (suite *APISuite) setupGin() {
	handler := NewRestfulHandler()
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()

	router := g.Group("/tinder_system/v1")

	router.POST("/persons", handler.AddSinglePersonAndMatch)
	router.DELETE("/persons/:name", handler.RemoveSinglePerson)
	router.GET("/persons", handler.QuerySinglePeople)
	ts := httptest.NewServer(g)
	suite.HTTPServer = ts
}

func (suite *APISuite) setupData() {

	person := models.SinglePerson{
		Name:         "test man",
		Height:       178,
		NumberOfDate: 10,
		Gender:       c.Male,
	}

	person2 := models.SinglePerson{
		Name:         "test feman",
		Height:       150,
		NumberOfDate: 3,
		Gender:       c.Female,
	}

	suite.testData = []*models.SinglePerson{&person, &person2}
}

func (suite *APISuite) TestAddSinglePersonAndMatch() {

	url := fmt.Sprintf("%s%s", suite.HTTPServer.URL, "/tinder_system/v1/persons")
	bs, _ := json.Marshal(suite.testData[0])
	addResp1, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		suite.Fail("HTTP POST request error:", err)
		return
	}

	// 檢查 HTTP 狀態碼是否為 200 OK
	if addResp1.StatusCode != http.StatusOK {
		suite.Fail("Expected status 200 OK, but got:", addResp1.Status)
		return
	}

	bs, _ = json.Marshal(suite.testData[1])
	addResp2, err := http.Post(url, "application/json", bytes.NewReader(bs))
	suite.Suite.Empty(err)

	defer addResp2.Body.Close()

	// 檢查 HTTP 狀態碼是否為 200 OK
	if addResp2.StatusCode != http.StatusOK {
		suite.Fail("Expected status 200 OK, but got:", addResp2.Status)
		return
	}

	addSinglePersonAndMatchResp := models.AddSinglePersonAndMatchResponse{}
	err = json.NewDecoder(addResp2.Body).Decode(&addSinglePersonAndMatchResp)
	suite.Suite.Empty(err)

	suite.testData[0].NumberOfDate--
	suite.testData[1].NumberOfDate--
	expect := models.AddSinglePersonAndMatchResponse{
		NewUser: suite.testData[1],
		Matches: []*models.SinglePerson{
			suite.testData[0],
		},
	}
	suite.Suite.EqualValues(expect, addSinglePersonAndMatchResp)

	suite.testQuerySinglePeople()
	suite.testRemoveSinglePerson()
}

// / TestQuerySinglePeople 由於會依賴於 add single person 的關係，所以將func name 改為小寫，並由 TestAddSinglePersonAndMatch 呼叫
func (suite *APISuite) testQuerySinglePeople() {
	uri := fmt.Sprintf("%s%s", suite.HTTPServer.URL, "/tinder_system/v1/persons")

	//test case 1
	queryParams := url.Values{}
	queryParams.Add("name", suite.testData[0].Name)
	queryParams.Add("n", "1")
	queryUrl := fmt.Sprintf("%s?%s", uri, queryParams.Encode())

	queryResp, err := http.Get(queryUrl)
	suite.Suite.Empty(err)
	defer queryResp.Body.Close()

	// 檢查 HTTP 狀態碼是否為 200 OK
	if queryResp.StatusCode != http.StatusOK {
		suite.Fail("Expected status 200 OK, but got:", queryResp.Status)
		return
	}

	querySinglePeopleResp := models.QuerySinglePeopleResponse{}
	err = json.NewDecoder(queryResp.Body).Decode(&querySinglePeopleResp)
	suite.Suite.Empty(err)

	suite.Suite.EqualValues(models.QuerySinglePeopleResponse{
		Matches: []*models.SinglePerson{suite.testData[0]},
	}, querySinglePeopleResp)

	//test case 2
	queryParams = url.Values{}
	queryParams.Add("heightGte", "150")
	queryParams.Add("n", "10")
	queryUrl = fmt.Sprintf("%s?%s", uri, queryParams.Encode())

	queryResp, err = http.Get(queryUrl)
	suite.Suite.Empty(err)
	defer queryResp.Body.Close()

	// 檢查 HTTP 狀態碼是否為 200 OK
	if queryResp.StatusCode != http.StatusOK {
		suite.Fail("Expected status 200 OK, but got:", queryResp.Status)
		return
	}

	querySinglePeopleResp = models.QuerySinglePeopleResponse{}
	err = json.NewDecoder(queryResp.Body).Decode(&querySinglePeopleResp)
	suite.Suite.Empty(err)

	suite.Suite.EqualValues(models.QuerySinglePeopleResponse{
		Matches: suite.testData,
	}, querySinglePeopleResp)
}

func (suite *APISuite) testRemoveSinglePerson() {
	url := fmt.Sprintf("%s%s/%s", suite.HTTPServer.URL, "/tinder_system/v1/persons", suite.testData[0].Name)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	suite.Suite.Empty(err)

	client := http.Client{}
	resp, err := client.Do(req)
	suite.Suite.Empty(err)
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼是否為 200 OK
	if resp.StatusCode != http.StatusOK {
		suite.Fail("Expected status 200 OK, but got:", resp.Status)
		return
	}

	commonResponse := models.CommonResponse{}
	err = json.NewDecoder(resp.Body).Decode(&commonResponse)
	suite.Suite.Empty(err)

	suite.Suite.EqualValues(models.CommonResponse{
		Result: "OK",
	}, commonResponse)
}
