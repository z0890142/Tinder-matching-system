package models

import "sync"

type SinglePerson struct {
	ID           int
	Name         string `json:"name"`
	Height       int    `json:"height"`
	Gender       string `json:"gender" example:"M or F"`
	NumberOfDate int    `json:"number_of_date"`
	Lock         sync.Mutex
}
