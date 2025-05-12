package geo

import (
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	AddressSuggestions(address string, limit int) ([]*suggest.AddressSuggestion, error)
	AddressSuggestionsValues(address string, limit int) ([]string, error)

	GetRegions(limit, offset int) ([]Region, error)
	/*
		AddRegion(name string) error // TODO: Implement this in the future.
	*/
	GetCitiesInRegion(regionId uuid.UUID, limit, offset int) ([]City, error)

	GetCities(limit, offset int) ([]City, error)
	/*
		AddCity(regionID, name string) error // TODO: Implement this in the future.
	*/
	GetBuildingsInCity(cityId uuid.UUID, limit, offset int) ([]Building, error)

	GetBuildings(limit, offset int) ([]Building, error)
	/*
		AddBuilding(cityID, street, houseNumber string) error // TODO: Implement this in the future.
	*/
	GetAddressesInBuilding(buildingId uuid.UUID, limit, offset int) ([]Address, error)

	AddFullAddress(regionName, cityName, street, houseNumber, apartmentNumber string) (*FullAddress, error)
}
