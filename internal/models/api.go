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
	Matches []*SinglePerson `json:"matches"`
}

type QueryFilter struct {
	Name      string
	N         int
	HeightGte int
	HeightLte int
	Gender    string
}

type CommonResponse struct {
	Result string
}
