package geo

import (
	"github.com/google/uuid"
)

type Region struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
	Lat  float64   `json:"lat" db:"lat"`
	Lon  float64   `json:"lon" db:"lon"`
}

type City struct {
	ID       uuid.UUID `json:"id" db:"id"`
	RegionID uuid.UUID `json:"region_id" db:"region_id"`
	Name     string    `json:"name" db:"name"`
	Lat      float64   `json:"lat" db:"lat"`
	Lon      float64   `json:"lon" db:"lon"`
}

type Building struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CityID      uuid.UUID `json:"city_id" db:"city_id"`
	Street      string    `json:"street" db:"street"`
	Additional  *string   `json:"additional" db:"additional"`
	Lat         float64   `json:"lat" db:"lat"`
	Lon         float64   `json:"lon" db:"lon"`
	HouseNumber string    `json:"house_number" db:"house_number"`
}

type Address struct {
	ID              uuid.UUID `json:"id" db:"id"`
	BuildingID      uuid.UUID `json:"building_id" db:"building_id"`
	ApartmentNumber string    `json:"apartment_number" db:"apartment_number"`
}

type FullAddress struct {
	RegionID        uuid.UUID `json:"region_id"`
	Region          string    `json:"region"`
	RegionLat       float64   `json:"region_lat"`
	RegionLon       float64   `json:"region_lon"`
	CityID          uuid.UUID `json:"city_id"`
	City            string    `json:"city"`
	CityLat         float64   `json:"city_lat"`
	CityLon         float64   `json:"city_lon"`
	BuildingID      uuid.UUID `json:"building_id"`
	Street          string    `json:"street"`
	StreetLat       float64   `json:"street_lat"`
	StreetLon       float64   `json:"street_lon"`
	HouseNumber     string    `json:"house_number"`
	AddressID       uuid.UUID `json:"address_id"`
	ApartmentNumber string    `json:"apartment_number"`
	FullAddress     string    `json:"full_address"`
}
