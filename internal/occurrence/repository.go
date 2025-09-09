package occurrence

import (
	"errors"
	"fmt"
)

type Repository interface {
	GetByID(id string) (*Occurrence, error)
	Save(oc Occurrence) error
	Remove(oc Occurrence) error
}

type memo struct {
	occurrenceTable map[string]*Occurrence
}

func NewMemoRepository() Repository {
	return &memo{
		occurrenceTable: map[string]*Occurrence{},
	}
}

func (m memo) Save(oc Occurrence) error {
	m.occurrenceTable[oc.ID] = &oc
	fmt.Println(&m.occurrenceTable)
	return nil
}

func (m memo) Remove(oc Occurrence) error {
	delete(m.occurrenceTable, oc.ID)
	return nil
}

func (m memo) GetByID(id string) (*Occurrence, error) {
	if cit, ok := m.occurrenceTable[id]; ok {
		return cit, nil
	}
	return nil, errors.New("invalid user id")
}
