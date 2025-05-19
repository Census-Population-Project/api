package geo

import (
	"context"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Config *config.Config

	DB     *database.DataBase
	Logger *logrus.Logger

	DDsAPI *suggest.Api

	CRUDGeo *CRUDGeo
}

func (s *Service) AddressSuggestions(address string, limit int) ([]*suggest.AddressSuggestion, error) {
	suggestions, err := s.DDsAPI.Address(context.Background(), &suggest.RequestParams{
		Query:    address,
		Count:    limit,
		Language: "RU",
	})
	if err != nil {
		s.Logger.Error("Error getting address suggestions: ", err)
		return nil, err
	}

	return suggestions, nil
}

func (s *Service) AddressSuggestionsValues(address string, limit int) ([]string, error) {
	suggestions, err := s.DDsAPI.Address(context.Background(), &suggest.RequestParams{
		Query:    address,
		Count:    limit,
		Language: "RU",
	})
	if err != nil {
		s.Logger.Error("Error getting address suggestions values: ", err)
		return nil, err
	}

	addressesValues := make([]string, len(suggestions))
	for i, suggestion := range suggestions {
		addressesValues[i] = suggestion.Value
	}

	return addressesValues, nil
}

func (s *Service) GetRegions(limit, offset int) ([]Region, error) {
	return s.CRUDGeo.SelectRegions(limit, offset)
}

func (s *Service) GetCitiesInRegion(regionId uuid.UUID, limit, offset int) ([]City, error) {
	return s.CRUDGeo.SelectCitiesInRegion(regionId, limit, offset)
}

func (s *Service) GetCities(limit, offset int) ([]City, error) {
	return s.CRUDGeo.SelectCities(limit, offset)
}

func (s *Service) GetBuildingsInCity(cityId uuid.UUID, limit, offset int) ([]Building, error) {
	return s.CRUDGeo.SelectBuildingsInCity(cityId, limit, offset)
}

func (s *Service) GetBuildings(limit, offset int) ([]Building, error) {
	return s.CRUDGeo.SelectBuildings(limit, offset)
}

func (s *Service) GetAddressesInBuilding(buildingId uuid.UUID, limit, offset int) ([]Address, error) {
	return s.CRUDGeo.SelectAddressesInBuilding(buildingId, limit, offset)
}

func (s *Service) AddFullAddress(
	regionName string, regionLat float64, regionLon float64,
	cityName string, cityLat float64, cityLon float64,
	street string, additional *string, streetLat float64, streetLon float64,
	houseNumber, apartmentNumber string,
) (*FullAddress, error) {
	return s.CRUDGeo.InsertFullAddress(
		regionName, regionLat, regionLon,
		cityName, cityLat, cityLon,
		street, additional, streetLat, streetLon,
		houseNumber, apartmentNumber,
	)
}

func NewService(cfg *config.Config, db *database.DataBase, logger *logrus.Logger, ddsApi *suggest.Api) *Service {
	return &Service{
		Config: cfg,

		DB:     db,
		Logger: logger,

		DDsAPI: ddsApi,

		CRUDGeo: NewGeoCRUD(db, logger),
	}
}
