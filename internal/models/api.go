package models

type SinglePersonMatchRequest struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

type AddSinglePersonAndMatchResponse struct {
	NewUser *SinglePerson   `json:"newUser"`
	Matches []*SinglePerson `json:"matches"`
}

type QuerySinglePeopleResponse struct {
	SinglePeople []*SinglePerson `json:"singlePeople"`
}

type QueryFilter struct {
	Name      string
	N         int
	HeightGte int
	HeightLte int
	Gender    string
}
