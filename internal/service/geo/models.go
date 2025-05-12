package geo

import "github.com/google/uuid"

type Region struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type City struct {
	ID       uuid.UUID `json:"id"`
	RegionID uuid.UUID `json:"region_id"`
	Name     string    `json:"name"`
}

type Building struct {
	ID          uuid.UUID `json:"id"`
	CityID      uuid.UUID `json:"city_id"`
	Street      string    `json:"street"`
	HouseNumber string    `json:"house_number"`
}

type Address struct {
	ID              uuid.UUID `json:"id"`
	BuildingID      uuid.UUID `json:"building_id"`
	ApartmentNumber string    `json:"apartment_number"`
}
