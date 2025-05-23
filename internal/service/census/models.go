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

type Language struct {
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
	PopulationCount int          `json:"population_count" db:"population_count"`
	Genders         []Gender     `json:"genders" db:"genders"`
	Regions         []geo.Region `json:"regions" db:"regions"`
	Cities          []geo.City   `json:"cities" db:"cities"`
}

type EventDataInLocation struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	StartDateTime   time.Time `json:"start_datetime" db:"start_datetime"`
	EndDateTime     time.Time `json:"end_datetime" db:"end_datetime"`
	PopulationCount int       `json:"population_count" db:"population_count"`
	Genders         []Gender  `json:"genders" db:"genders"`
}

type Events struct {
	Events []Event `json:"events"`
	Total  int64   `json:"total"`
}

type EventStatistics struct {
	// General information
	TotalPopulation        int     `json:"total_population"`
	TotalHouseholds        int     `json:"total_households"`
	AvgPersonsPerHousehold float64 `json:"avg_persons_per_household"`

	// Population structure
	GenderDistribution []Gender `json:"gender_distribution"`
	AverageAge         float64  `json:"average_age"`
	ChildrenCount      int      `json:"children_count"`
	ElderlyCount       int      `json:"elderly_count"`

	// Education and employment
	EducationDistribution  map[string]int `json:"education_distribution"`
	EmploymentDistribution map[string]int `json:"employment_distribution"`

	// Languages and citizenship
	PercentSpeaksRussian float64    `json:"percent_speaks_russian"`
	DualCitizenshipCount int        `json:"dual_citizenship_count"`
	TopOtherLanguages    []Language `json:"top_other_languages"`

	IncomeSourcesDistribution map[string]int `json:"income_sources_distribution"`
}
