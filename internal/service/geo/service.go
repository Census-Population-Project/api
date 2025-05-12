package geo

type ServiceInterface interface {
	GetRegions(limit, offset int) ([]Region, error)
	AddRegion(name string) error
	GetCitiesInRegion(regionID string, limit, offset int) ([]City, error)

	GetCities(limit, offset int) ([]City, error)
	AddCity(regionID, name string) error
	GetBuildingsInCity(cityID string, limit, offset int) ([]Building, error)

	GetBuildings(limit, offset int) ([]Building, error)
	AddBuilding(cityID, street, houseNumber string) error
	GetAddressesInBuilding(buildingID string, limit, offset int) ([]Address, error)
}
