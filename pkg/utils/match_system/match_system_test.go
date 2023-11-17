package matchsystem

import (
	"testing"
	"tinderMatchingSystem/internal/c"
	"tinderMatchingSystem/internal/models"

	"github.com/mohae/deepcopy"

	"github.com/stretchr/testify/assert"
)

var heightIndex = map[int][]int{
	150: []int{1, 2, 3, 4, 5, 6},
	160: []int{7},
	180: []int{8},
	200: []int{9},
}

func TestPersonValidation(t *testing.T) {
	tests := []struct {
		name    string
		person  *models.SinglePerson
		wantErr bool
	}{
		{
			"Success Test",
			&models.SinglePerson{
				Name:         "success test",
				Height:       180,
				Gender:       c.Female,
				NumberOfDate: 1,
			},
			false,
		},
		{
			"Miss Name",
			&models.SinglePerson{
				Height:       180,
				Gender:       c.Female,
				NumberOfDate: 1,
			},
			true,
		},
		{
			"Miss Heigth",
			&models.SinglePerson{
				Name:         "miss heigth",
				Gender:       c.Female,
				NumberOfDate: 1,
			},
			true,
		},
		{
			"Error Gender",
			&models.SinglePerson{
				Name:         "miss heigth",
				Height:       110,
				Gender:       "test",
				NumberOfDate: 1,
			},
			true,
		},
		{
			"miss NumberOfDate",
			&models.SinglePerson{
				Name:   "miss heigth",
				Height: 110,
				Gender: c.Male,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePerson(tt.person); (err != nil) != tt.wantErr {
				t.Errorf("personValidation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindHeightIdx(t *testing.T) {

	tests := []struct {
		name      string
		heightGte int
		heightLte int
		expect    *[]int
	}{
		{
			"heigth >=120",
			120,
			0,
			&[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			"heigth <=160",
			0,
			160,
			&[]int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"heigth >=160 && height <= 200",
			160,
			200,
			&[]int{7, 8, 9},
		},
		{
			"heigth >=201",
			201,
			0,
			&[]int{9},
		},
	}
	for _, tt := range tests {
		result := []int{}
		FindHeightIdx(heightIndex, tt.heightGte, tt.heightLte, &result)
		assert.EqualValues(t, tt.expect, &result, tt.name)
	}
}

func TestDeleteHeightIndex(t *testing.T) {
	tests := []struct {
		name        string
		heightLevel int
		target      int
		expect      map[int][]int
	}{

		{
			"Success",
			150,
			4,
			map[int][]int{
				150: []int{1, 2, 3, 5, 6},
				160: []int{7},
				180: []int{8},
				200: []int{9},
			},
		},
		{
			"not delete",
			150,
			7,
			map[int][]int{
				150: []int{1, 2, 3, 4, 5, 6},
				160: []int{7},
				180: []int{8},
				200: []int{9},
			},
		},
		{
			"delete all slice",
			200,
			9,
			map[int][]int{
				150: []int{1, 2, 3, 4, 5, 6},
				160: []int{7},
				180: []int{8},
			},
		},
	}
	for _, tt := range tests {
		indexMap := deepcopy.Copy(heightIndex).(map[int][]int)

		DeleteHeightIndex(indexMap, tt.target, tt.heightLevel)
		assert.EqualValues(t, tt.expect, indexMap, tt.name)
	}
}
