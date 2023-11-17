package matchesystem

import (
	"sort"
	"testing"
	"tinderMatchingSystem/internal/c"
	"tinderMatchingSystem/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewMatcheSystem(t *testing.T) {
	assert.NotEmpty(t, NewMatcheSystem())
}

func TestRegisterSinglePerson(t *testing.T) {
	matcheSystem := NewMatcheSystem()
	index := matcheSystem.GetHightIndex(c.Male)
	nameIndex := matcheSystem.GetNameIndex()

	person := &models.SinglePerson{
		Name:         "person1",
		Height:       170,
		Gender:       c.Male,
		NumberOfDate: 10,
	}

	err := matcheSystem.RegisterSinglePerson(person)
	assert.Empty(t, err)
	assert.Equal(t, []int{0}, index[170])
	assert.NotEmpty(t, nameIndex[person.Name])

	//register same person name
	err = matcheSystem.RegisterSinglePerson(person)
	assert.NotEmpty(t, err)
	assert.Equal(t, []int{0}, index[170])

	person2 := &models.SinglePerson{
		Name:         "person2",
		Height:       151,
		Gender:       c.Male,
		NumberOfDate: 10,
	}

	err = matcheSystem.RegisterSinglePerson(person2)
	assert.Empty(t, err)
	assert.Equal(t, []int{1}, index[150])

}

func TestRemoveSinglePerson(t *testing.T) {
	matcheSystem := NewMatcheSystem()
	index := matcheSystem.GetHightIndex(c.Female)
	nameIndex := matcheSystem.GetNameIndex()

	person := &models.SinglePerson{
		Name:         "person1",
		Height:       150,
		Gender:       c.Female,
		NumberOfDate: 1,
	}
	err := matcheSystem.RegisterSinglePerson(person)
	assert.Empty(t, err)
	assert.Equal(t, []int{0}, index[150])

	person2 := &models.SinglePerson{
		Name:         "person2",
		Height:       155,
		Gender:       c.Female,
		NumberOfDate: 11,
	}
	err = matcheSystem.RegisterSinglePerson(person2)
	assert.Empty(t, err)
	assert.Equal(t, []int{0, 1}, index[150])

	err = matcheSystem.RemoveSinglePerson("person1")
	assert.Empty(t, err)
	assert.Equal(t, []int{1}, index[150])
	assert.Empty(t, nameIndex["person1"])

}

func TestQuerySinglePerson(t *testing.T) {
	matcheSystem := NewMatcheSystem()
	persons := []*models.SinglePerson{
		&models.SinglePerson{
			Name:         "male1",
			Height:       120,
			Gender:       c.Male,
			NumberOfDate: 10,
		},
		&models.SinglePerson{
			Name:         "male2",
			Height:       157,
			Gender:       c.Male,
			NumberOfDate: 1,
		},
		&models.SinglePerson{
			Name:         "male3",
			Height:       165,
			Gender:       c.Male,
			NumberOfDate: 5,
		},
		&models.SinglePerson{
			Name:         "female1",
			Height:       167,
			Gender:       c.Female,
			NumberOfDate: 10,
		},
		&models.SinglePerson{
			Name:         "female2",
			Height:       171,
			Gender:       c.Female,
			NumberOfDate: 5,
		}, &models.SinglePerson{
			Name:         "female3",
			Height:       180,
			Gender:       c.Female,
			NumberOfDate: 1,
		},
	}

	for _, p := range persons {
		assert.Empty(t, matcheSystem.RegisterSinglePerson(p))
	}

	result, err := matcheSystem.QuerySinglePerson(&models.QueryFilter{
		Name: "male3",
	})
	assert.NotEmpty(t, err)

	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		Name: "male3",
		N:    1,
	})
	assert.Empty(t, err)
	assert.Equal(t, []*models.SinglePerson{persons[2]}, result)

	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		Gender: c.Female,
		N:      10,
	})
	assert.Empty(t, err)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Height < result[j].Height
	})
	assert.Equal(t, persons[3:], result)

	//query by hight range
	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		N:         10,
		HeightGte: 155,
		HeightLte: 170,
	})
	assert.Empty(t, err)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Height < result[j].Height
	})
	assert.Equal(t, persons[1:4], result)

	//query by hight and gender range
	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		N:         10,
		HeightGte: 155,
		HeightLte: 175,
		Gender:    c.Female,
	})
	assert.Empty(t, err)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Height < result[j].Height
	})
	assert.Equal(t, persons[3:5], result)

	//query by hight > 150
	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		N:         10,
		HeightGte: 150,
	})
	assert.Empty(t, err)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Height < result[j].Height
	})
	assert.Equal(t, persons[1:], result)

	//query by hight > 150 and gender
	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		N:         10,
		HeightGte: 150,
		Gender:    c.Male,
	})
	assert.Empty(t, err)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Height < result[j].Height
	})
	assert.Equal(t, persons[1:3], result)

	// query by hight < 150
	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		N:         10,
		HeightLte: 150,
	})
	assert.Empty(t, err)
	assert.Equal(t, []*models.SinglePerson{persons[0]}, result)

	//query by hight < 150 and gender
	result, err = matcheSystem.QuerySinglePerson(&models.QueryFilter{
		N:         10,
		HeightLte: 150,
		Gender:    c.Female,
	})
	assert.Empty(t, err)
	assert.Equal(t, 0, len(result))
}
