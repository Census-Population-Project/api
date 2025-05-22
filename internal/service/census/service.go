package census

import "github.com/google/uuid"

type ServiceInterface interface {
	GetEvents(limit, offset int) ([]Event, *int64, error)
	GetEventInfoByID(id uuid.UUID) (*EventInfo, error)
	GetEventInfoByLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventDataInLocation, error)
}
