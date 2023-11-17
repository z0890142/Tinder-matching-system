package matchesystem

import (
	"fmt"
	"sync"
	"tinderMatchingSystem/internal/c"
	"tinderMatchingSystem/internal/models"
	matchSystemUtils "tinderMatchingSystem/pkg/utils/match_system"
)

type MatcheSystem interface {
	RegisterSinglePerson(person *models.SinglePerson) error
	RemoveSinglePerson(personName string) error
	MatchingPerson(person *models.SinglePerson, numberOfDate int) []*models.SinglePerson
	QuerySinglePerson(filter *models.QueryFilter) ([]*models.SinglePerson, error)

	GetHightIndex(gender string) map[int][]int
	GetNameIndex() map[string]*models.SinglePerson
}

type matcheSystem struct {
	mu sync.RWMutex

	singlePersons     map[int]*models.SinglePerson
	nameIndex         map[string]*models.SinglePerson
	maleHeightIndex   map[int][]int
	femaleHeightIndex map[int][]int

	nextID int
}

func NewMatcheSystem() MatcheSystem {

	return &matcheSystem{
		singlePersons:     make(map[int]*models.SinglePerson),
		nameIndex:         make(map[string]*models.SinglePerson),
		maleHeightIndex:   make(map[int][]int),
		femaleHeightIndex: make(map[int][]int),
	}
}

func (s *matcheSystem) RegisterSinglePerson(person *models.SinglePerson) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := matchSystemUtils.ValidatePerson(person); err != nil {
		return fmt.Errorf("RegisterSinglePerson : %s", err.Error())
	}

	person.ID = s.nextID

	if _, ok := s.nameIndex[person.Name]; !ok {
		s.nameIndex[person.Name] = person
	} else {
		return fmt.Errorf("RegisterSinglePerson : person alrealdy exist")
	}

	heightLevel := person.Height - (person.Height % 10)
	s.singlePersons[person.ID] = person

	if person.Gender == c.Male {
		s.maleHeightIndex[heightLevel] = append(s.maleHeightIndex[heightLevel], person.ID)
	} else if person.Gender == c.Female {
		s.femaleHeightIndex[heightLevel] = append(s.femaleHeightIndex[heightLevel], person.ID)
	}

	s.nextID++

	return nil
}

func (s *matcheSystem) RemoveSinglePerson(personName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if personName == "" {
		return fmt.Errorf("RegisterSinglePerson : person name is null")
	}
	person, ok := s.nameIndex[personName]
	if !ok {
		return fmt.Errorf("RegisterSinglePerson : person not exist")
	}

	heightLevel := person.Height - (person.Height % 10)

	if person.Gender == c.Male {
		matchSystemUtils.DeleteHeightIndex(s.maleHeightIndex, person.ID, heightLevel)
	} else {
		matchSystemUtils.DeleteHeightIndex(s.femaleHeightIndex, person.ID, heightLevel)
	}

	delete(s.nameIndex, personName)
	delete(s.singlePersons, person.ID)

	return nil
}

func (s *matcheSystem) QuerySinglePerson(filter *models.QueryFilter) ([]*models.SinglePerson, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if filter.N == 0 {
		return nil, fmt.Errorf("QuerySinglePerson : n can't be 0")
	}
	result := make([]*models.SinglePerson, 0)

	if filter.Name != "" {
		person, ok := s.nameIndex[filter.Name]
		if !ok {
			return nil, fmt.Errorf("QuerySinglePerson : person not found")
		}
		result = append(result, person)
		return result, nil
	}

	heightIdxs := make([]int, 0)
	if filter.Gender == "" {
		matchSystemUtils.FindHeightIdx(s.maleHeightIndex, filter.HeightGte, filter.HeightLte, &heightIdxs)
		matchSystemUtils.FindHeightIdx(s.femaleHeightIndex, filter.HeightGte, filter.HeightLte, &heightIdxs)
	} else if filter.Gender == c.Male {
		matchSystemUtils.FindHeightIdx(s.maleHeightIndex, filter.HeightGte, filter.HeightLte, &heightIdxs)
	} else {
		matchSystemUtils.FindHeightIdx(s.femaleHeightIndex, filter.HeightGte, filter.HeightLte, &heightIdxs)
	}

	for _, idx := range heightIdxs {
		if filter.N == 0 {
			break
		}
		target := s.singlePersons[idx]
		if (filter.HeightLte != 0 && target.Height > filter.HeightLte) || (filter.HeightGte != 0 && target.Height < filter.HeightGte) {
			continue
		}
		result = append(result, target)
		filter.N--
	}
	return result, nil
}

func (s *matcheSystem) MatchingPerson(person *models.SinglePerson, numberOfDate int) []*models.SinglePerson {
	heightLevel := person.Height - (person.Height % 10)
	matches := make([]*models.SinglePerson, 0)
	heightIdxs := make([]int, 0)

	if person.Gender == c.Female {
		matchSystemUtils.FindHeightIdx(s.maleHeightIndex, heightLevel, 0, &heightIdxs)
	} else {
		matchSystemUtils.FindHeightIdx(s.femaleHeightIndex, 0, heightLevel, &heightIdxs)
	}

	for _, idx := range heightIdxs {
		if person.NumberOfDate == 0 {
			break
		}

		target := s.singlePersons[idx]
		target.Lock.Lock()

		if target.NumberOfDate == 0 {
			s.RemoveSinglePerson(target.Name)
			target.Lock.Unlock()
			continue
		}

		if (person.Gender == c.Male && target.Height >= person.Height) || (person.Gender == c.Female && target.Height <= person.Height) {
			target.Lock.Unlock()
			continue
		}

		matches = append(matches, target)
		person.NumberOfDate--
		target.NumberOfDate--
		target.Lock.Unlock()
	}
	return matches
}

func (s *matcheSystem) GetHightIndex(gender string) map[int][]int {
	if gender == c.Male {
		return s.maleHeightIndex
	} else if gender == c.Female {
		return s.femaleHeightIndex
	}
	return nil
}

func (s *matcheSystem) GetNameIndex() map[string]*models.SinglePerson {
	return s.nameIndex
}
