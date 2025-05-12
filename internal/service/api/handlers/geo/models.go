package geo

import "github.com/google/uuid"

type CreateRegionRequest struct {
	Name string `json:"name"`
}

type CreateCityRequest struct {
	Name     string    `json:"name"`
	RegionID uuid.UUID `json:"region_id"`
}

type CreateBuildingRequest struct {
	CityID      uuid.UUID `json:"city_id"`
	Street      string    `json:"street"`
	HouseNumber string    `json:"house_number"`
}
