package citizen

import "errors"

type Repository interface {
	GetByID(id string) (*Citizen, error)
	GetByEmail(email string) (*Citizen, error)
	Save(cit *Citizen) error
}

type memo struct {
	citizenTable map[string]*Citizen
}

func NewMemoRepository() Repository {
	return &memo{
		citizenTable: map[string]*Citizen{},
	}
}

func (m memo) GetByID(id string) (*Citizen, error) {
	if cit, ok := m.citizenTable[id]; ok {
		return cit, nil
	}
	return nil, errors.New("invalid user id")
}

func (m memo) GetByEmail(email string) (*Citizen, error) {
	if cit, ok := m.citizenTable[email]; ok {
		return cit, nil
	}
	return nil, errors.New("invalid user id")
}

func (m *memo) Save(cit *Citizen) error {
	m.citizenTable[cit.ID] = cit
	m.citizenTable[cit.Email] = cit
	return nil
}
