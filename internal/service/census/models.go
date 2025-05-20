package census

import (
	"time"

	"github.com/Census-Population-Project/API/internal/service/geo"

	"github.com/google/uuid"
)

type Gender struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type AgeGroup struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type Event struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	StartDateTime time.Time `json:"start_datetime" db:"start_datetime"`
	EndDateTime   time.Time `json:"end_datetime" db:"end_datetime"`
}

type EventInfo struct {
	ID              uuid.UUID    `json:"id" db:"id"`
	Name            string       `json:"name" db:"name"`
	StartDateTime   time.Time    `json:"start_datetime" db:"start_datetime"`
	EndDateTime     time.Time    `json:"end_datetime" db:"end_datetime"`
	RegionsCount    int          `json:"regions_count" db:"regions_count"`
	CitiesCount     int          `json:"cities_count" db:"cities_count"`
	PopulationCount int          `json:"population" db:"population"`
	Genders         []Gender     `json:"genders" db:"genders"`
	Regions         []geo.Region `json:"regions" db:"regions"`
	Cities          []geo.City   `json:"cities" db:"cities"`
}

type EventDataInLocation struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	StartDateTime   time.Time `json:"start_datetime" db:"start_datetime"`
	EndDateTime     time.Time `json:"end_datetime" db:"end_datetime"`
	PopulationCount int       `json:"population" db:"population"`
	Genders         []Gender  `json:"genders" db:"genders"`
}
