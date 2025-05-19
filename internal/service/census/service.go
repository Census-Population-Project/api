package census

import "github.com/google/uuid"

type ServiceInterface interface {
	GetEvents(limit, offset int) ([]Event, error)
	GetEventInfoByID(id uuid.UUID) (*EventInfo, error)
	GetEventInfoByLocationIDs(regionId uuid.UUID, cityId uuid.UUID) (*EventInfo, error)
}
