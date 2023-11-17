package matchsystem

import (
	"fmt"
	"tinderMatchingSystem/internal/c"
	"tinderMatchingSystem/internal/models"
)

func FindHeightIdx(heightIndex map[int][]int, heightGte, heightLte int, heightIdx *[]int) {

	heighGteLevel := heightGte - (heightGte % 10)
	heighLteLevel := heightLte - (heightLte % 10)

	if heightLte != 0 && heightGte != 0 {
		for k, v := range heightIndex {
			if k >= heighGteLevel && k <= heighLteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else if heightLte == 0 && heightGte != 0 {
		for k, v := range heightIndex {
			if k >= heighGteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else if heightLte != 0 && heightGte == 0 {

		for k, v := range heightIndex {
			if k <= heighLteLevel {
				*heightIdx = append(*heightIdx, v...)
			}
		}
	} else {
		for _, v := range heightIndex {
			*heightIdx = append(*heightIdx, v...)

		}
	}
	return
}

func DeleteHeightIndex(index map[int][]int, targetId, heightLevel int) error {
	if indices, ok := index[heightLevel]; ok {
		for i, idx := range indices {
			if idx == targetId {
				index[heightLevel] = append(indices[:i], indices[i+1:]...)
				break
			}
		}
		if len(index[heightLevel]) == 0 {
			delete(index, heightLevel)
		}
	}

	return nil
}

func ValidatePerson(person *models.SinglePerson) error {
	if person.Gender != c.Male && person.Gender != c.Female {
		return fmt.Errorf("personValidation : gender can't empty")
	}

	if person.Height <= 0 {
		return fmt.Errorf("personValidation : height validate error")
	}

	if person.Name == "" {
		return fmt.Errorf("personValidation : name can't empty")
	}

	if person.NumberOfDate <= 0 {
		return fmt.Errorf("personValidation : NumberOfDate need over 0")
	}
	return nil
}
