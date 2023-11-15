package matchesystem

import (
	"fmt"
	"sync"
	"tinderMatchingSystem/internal/models"
)

type MatcheSystem interface {
	RegisterSinglePerson(person *models.SinglePerson) error
	RemoveSinglePerson(personName string) error
	MatchingFmalePerson(person *models.SinglePerson, numberOfDate int) []*models.SinglePerson
	MatchingMalePerson(person *models.SinglePerson, numberOfDate int) []*models.SinglePerson
	QuerySinglePerson(filter models.QueryFilter) ([]*models.SinglePerson, error)
}

type matcheSystem struct {
	mu            sync.RWMutex
	singlePersons sync.Map

	nameIndex        sync.Map
	maleHeightIndex  map[int][]int
	fmaleHeightIndex map[int][]int

	nextID int
}

func NewMatcheSystem() MatcheSystem {

	return &matcheSystem{
		maleHeightIndex:  make(map[int][]int),
		fmaleHeightIndex: make(map[int][]int),
	}
}

func (s *matcheSystem) RegisterSinglePerson(person *models.SinglePerson) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	person.ID = s.nextID

	if _, load := s.nameIndex.LoadOrStore(person.Name, &person); load {
		fmt.Println("RegisterSinglePerson : person alrealdy exist")
	}

	heightLevel := person.Height - (person.Height % 10)
	s.singlePersons.Store(person.ID, person)

	if person.Gender == "M" {
		s.maleHeightIndex[heightLevel] = append(s.maleHeightIndex[heightLevel], person.ID)
	} else if person.Gender == "F" {
		s.fmaleHeightIndex[heightLevel] = append(s.fmaleHeightIndex[heightLevel], person.ID)
	}

	s.nextID++

	return nil
}

func (s *matcheSystem) RemoveSinglePerson(personName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.nameIndex.Load(personName)
	if !ok {
		fmt.Println("RegisterSinglePerson : person not exist")
	}

	person := p.(*models.SinglePerson)
	heightLevel := person.Height - (person.Height % 10)

	if person.Gender == "M" {
		// delete height index
		if indices, ok := s.maleHeightIndex[heightLevel]; ok {
			for i, idx := range indices {
				if idx == person.ID {
					s.maleHeightIndex[heightLevel] = append(indices[:i], indices[i+1:]...)
					break
				}
			}
			if len(s.maleHeightIndex[heightLevel]) == 0 {
				delete(s.maleHeightIndex, heightLevel)
			}
		}

	} else {
		if indices, ok := s.fmaleHeightIndex[heightLevel]; ok {
			for i, idx := range indices {
				if idx == person.ID {
					s.fmaleHeightIndex[heightLevel] = append(indices[:i], indices[i+1:]...)
					break
				}
			}
			if len(s.fmaleHeightIndex[heightLevel]) == 0 {
				delete(s.fmaleHeightIndex, heightLevel)
			}
		}
	}

	s.nameIndex.Delete(personName)

	return nil
}

func (s *matcheSystem) QuerySinglePerson(filter models.QueryFilter) ([]*models.SinglePerson, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.SinglePerson

	if filter.Name != "" {
		p, ok := s.nameIndex.Load(filter.Name)
		if !ok {
			return nil, fmt.Errorf("QuerySinglePerson : person not found")
		}
		result = append(result, p.(*models.SinglePerson))
		return result, nil
	}

	heightIdxs := make([]int, 0)
	if filter.Gender == "" {
		s.findFmaleHeightIdx(filter.HeightGte, filter.HeightLte, &heightIdxs)
		s.findMaleHeightIdx(filter.HeightGte, filter.HeightLte, &heightIdxs)
	} else if filter.Gender == "M" {
		s.findMaleHeightIdx(filter.HeightGte, filter.HeightLte, &heightIdxs)
	} else {
		s.findFmaleHeightIdx(filter.HeightGte, filter.HeightLte, &heightIdxs)

	}

	for _, idx := range heightIdxs {
		if filter.N == 0 {
			break
		}
		t, _ := s.singlePersons.Load(idx)
		target := t.(*models.SinglePerson)
		if target.Height > filter.HeightLte || target.Height < filter.HeightGte {
			continue
		}
		result = append(result, target)
		filter.N--
	}
	return result, nil
}

func (s *matcheSystem) MatchingMalePerson(person *models.SinglePerson, numberOfDate int) []*models.SinglePerson {

	matches := make([]*models.SinglePerson, 0)
	heightLevel := person.Height - (person.Height % 10)
	heightIdxs := make([]int, 0)
	s.findMaleHeightIdx(heightLevel, 0, &heightIdxs)

	for _, idx := range heightIdxs {
		if person.NumberOfDate == 0 {
			break
		}

		t, _ := s.singlePersons.Load(idx)
		target := t.(*models.SinglePerson)
		target.Lock.Lock()
		if target.NumberOfDate == 0 {
			s.mu.Lock()
			s.singlePersons.Delete(idx)
			s.mu.Unlock()
			target.Lock.Unlock()
			continue
		}
		if target.Height < person.Height {
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

func (s *matcheSystem) MatchingFmalePerson(person *models.SinglePerson, numberOfDate int) []*models.SinglePerson {

	matches := make([]*models.SinglePerson, 0)
	heightLevel := person.Height - (person.Height % 10)
	heightIdxs := make([]int, 0)
	s.findFmaleHeightIdx(heightLevel, 0, &heightIdxs)

	for _, idx := range heightIdxs {
		if person.NumberOfDate == 0 {
			break
		}

		t, _ := s.singlePersons.Load(idx)
		target := t.(*models.SinglePerson)
		target.Lock.Lock()
		if target.NumberOfDate == 0 {
			s.mu.Lock()
			s.singlePersons.Delete(idx)
			s.mu.Unlock()
			target.Lock.Unlock()
			continue
		}
		if target.Height > person.Height {
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

func (s *matcheSystem) findMaleHeightIdx(heightGte, heightLte int, heightIdx *[]int) {
	heighGteLevel := heightGte - (heightGte % 10)
	heighLteLevel := heightLte - (heightLte % 10)

	if heightLte != 0 && heightGte != 0 {
		for k, v := range s.maleHeightIndex {
			if k >= heighGteLevel && k <= heighLteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else if heightLte == 0 && heightGte != 0 {
		for k, v := range s.maleHeightIndex {
			if k >= heighGteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else if heightLte != 0 && heightGte == 0 {

		for k, v := range s.maleHeightIndex {
			if k <= heighLteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	}
	return
}

func (s *matcheSystem) findFmaleHeightIdx(heightGte, heightLte int, heightIdx *[]int) {
	heighGteLevel := heightGte - (heightGte % 10)
	heighLteLevel := heightLte - (heightLte % 10)

	if heightLte != 0 && heightGte != 0 {
		for k, v := range s.fmaleHeightIndex {
			if k >= heighGteLevel && k <= heighLteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else if heightLte == 0 && heightGte != 0 {
		for k, v := range s.fmaleHeightIndex {
			if k >= heighGteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else if heightLte != 0 && heightGte == 0 {

		for k, v := range s.fmaleHeightIndex {
			if k <= heighLteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	}
	return
}
