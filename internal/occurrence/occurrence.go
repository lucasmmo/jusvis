package occurrence

import (
	"jusvis/internal/citizen"

	"github.com/google/uuid"
)

type Type string

const (
	HOLE      Type = "HOLE"
	GARBAGE   Type = "GARBAGE"
	VANDALISM Type = "VANDALISM"
	THEFT     Type = "THEFT"
	UNKNOWN   Type = "UNKNOWN"
)

var ValidTypes = map[Type]bool{
	"HOLE":      true,
	"GARBAGE":   true,
	"VANDALISM": true,
	"THEFT":     true,
}

func getType(t string) Type {
	if _, ok := ValidTypes[Type(t)]; ok {
		return Type(t)
	}
	return UNKNOWN
}

type Occurrence struct {
	ID        string
	Type      Type
	RelatedBy string
}

type createCommand struct {
	occurrenceRepository Repository
	citizenRepository    citizen.Repository
}

func NewCreateCommand(
	occurrenceRepository Repository,
	citizenRepository citizen.Repository,
) *createCommand {
	return &createCommand{
		occurrenceRepository: occurrenceRepository,
		citizenRepository:    citizenRepository,
	}
}

func (c *createCommand) Do(ocType, relatedBy string) error {
	cit, err := c.citizenRepository.GetByID(relatedBy)
	if err != nil {
		return err
	}
	oc := Occurrence{
		ID:        uuid.NewString(),
		Type:      getType(ocType),
		RelatedBy: cit.ID,
	}

	if err := c.occurrenceRepository.Save(oc); err != nil {
		return err
	}
	return nil
}

func (c *createCommand) Undo(id string) error {
	oc, err := c.occurrenceRepository.GetByID(id)
	if err != nil {
		return err
	}
	if err := c.occurrenceRepository.Remove(*oc); err != nil {
		return err
	}
	return nil
}
