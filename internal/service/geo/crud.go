package geo

import (
	"context"
	"errors"

	"github.com/Census-Population-Project/API/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type CRUDInterface interface {
	SelectRegions(limit, offset int) ([]Region, error)
	InsertRegion(name string) (*uuid.UUID, error) // TODO: Implement this in the future.
	SelectCitiesInRegion(regionId uuid.UUID, limit, offset int) ([]City, error)

	SelectCities(limit, offset int) ([]City, error)
	InsertCity(regionId uuid.UUID, name string) (*uuid.UUID, error) // TODO: Implement this in the future.
	SelectBuildingsInCity(cityId uuid.UUID, limit, offset int) ([]Building, error)

	SelectBuildings(limit, offset int) ([]Building, error)
	InsertBuilding(cityId uuid.UUID, street, houseNumber string) (*uuid.UUID, error) // TODO: Implement this in the future.
	SelectAddressesInBuilding(buildingId uuid.UUID, limit, offset int) ([]Address, error)

	InsertFullAddress(
		regionName string, regionLat float64, regionLon float64,
		cityName string, cityLat float64, cityLon float64,
		street string, additional *string, streetLat float64, streetLon float64,
		houseNumber, apartmentNumber string,
	) (*FullAddress, error)
}

type CRUDGeo struct {
	DataBase *database.DataBase
	Logger   *logrus.Logger
}

func (s *CRUDGeo) SelectRegions(limit, offset int) ([]Region, error) {
	query := `SELECT id, name, lat, lon FROM geo.regions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to select regions: ", err)
		return nil, err
	}
	defer rows.Close()

	regions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Region])
	if err != nil {
		s.Logger.Error("Failed to collect regions: ", err)
		return nil, err
	}

	return regions, nil
}

func (s *CRUDGeo) InsertRegion(name string) (*uuid.UUID, error) {
	return nil, nil
}

func (s *CRUDGeo) SelectCitiesInRegion(regionId uuid.UUID, limit, offset int) ([]City, error) {
	query := `SELECT id, region_id, name, lat, lon FROM geo.cities WHERE region_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, regionId, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to select cities in region: ", err)
		return nil, err
	}
	defer rows.Close()

	cities, err := pgx.CollectRows(rows, pgx.RowToStructByName[City])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewRegionNotFoundError()
		}
		s.Logger.Error("Failed to collect cities: ", err)
		return nil, err
	}

	return cities, nil
}

func (s *CRUDGeo) SelectCities(limit, offset int) ([]City, error) {
	query := `SELECT id, region_id, name, lat, lon FROM geo.cities ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to select cities: ", err)
		return nil, err
	}
	defer rows.Close()

	cities, err := pgx.CollectRows(rows, pgx.RowToStructByName[City])
	if err != nil {
		s.Logger.Error("Failed to collect cities: ", err)
		return nil, err
	}

	return cities, nil
}

func (s *CRUDGeo) InsertCity(regionId uuid.UUID, name string) (*uuid.UUID, error) {
	return nil, nil
}

func (s *CRUDGeo) SelectBuildingsInCity(cityId uuid.UUID, limit, offset int) ([]Building, error) {
	query := `SELECT id, city_id, street, additional, house_number, lat, lon FROM geo.buildings WHERE city_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, cityId, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to select buildings in city: ", err)
		return nil, err
	}
	defer rows.Close()

	buildings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Building])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewCityNotFoundError()
		}
		s.Logger.Error("Failed to collect buildings: ", err)
		return nil, err
	}

	return buildings, nil
}

func (s *CRUDGeo) SelectBuildings(limit, offset int) ([]Building, error) {
	query := `SELECT id, city_id, street, additional, house_number, lat, lon FROM geo.buildings ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to select buildings: ", err)
		return nil, err
	}
	defer rows.Close()

	buildings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Building])
	if err != nil {
		s.Logger.Error("Failed to collect buildings: ", err)
		return nil, err
	}

	return buildings, nil
}

func (s *CRUDGeo) InsertBuilding(cityId uuid.UUID, street, houseNumber string) (*uuid.UUID, error) {
	return nil, nil
}

func (s *CRUDGeo) SelectAddressesInBuilding(buildingId uuid.UUID, limit, offset int) ([]Address, error) {
	query := `SELECT id, building_id, apartment_number FROM geo.addresses WHERE building_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := s.DataBase.DBPool.Query(context.Background(), query, buildingId, limit, offset)
	if err != nil {
		s.Logger.Error("Failed to select addresses in building: ", err)
		return nil, err
	}
	defer rows.Close()

	addresses, err := pgx.CollectRows(rows, pgx.RowToStructByName[Address])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewBuildingNotFoundError()
		}
		s.Logger.Error("Failed to collect addresses: ", err)
		return nil, err
	}

	return addresses, nil
}

func (s *CRUDGeo) InsertFullAddress(
	regionName string, regionLat float64, regionLon float64,
	cityName string, cityLat float64, cityLon float64,
	street string, additional *string, streetLat float64, streetLon float64,
	houseNumber, apartmentNumber string,
) (*FullAddress, error) {
	ctx := context.Background()
	tx, err := s.DataBase.DBPool.Begin(ctx)
	if err != nil {
		s.Logger.Error("Failed to begin transaction: ", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	regionQuery := `INSERT INTO geo.regions (name, lat, lon)
	VALUES ($1, $2, $3)
	ON CONFLICT (name) DO NOTHING
	RETURNING id, name, lat, lon`
	row, err := tx.Query(ctx, regionQuery, regionName, regionLat, regionLon)
	if err != nil {
		s.Logger.Error("Failed to insert/select region: ", err)
		return nil, err
	}
	defer row.Close()

	region, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Region])
	if err != nil {
		s.Logger.Error("Failed to collect region: ", err)
		return nil, err
	}

	cityQuery := `INSERT INTO geo.cities (name, region_id, lat, lon)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (name) DO NOTHING
	RETURNING id, region_id, name, lat, lon`
	row, err = tx.Query(ctx, cityQuery, cityName, region.ID, cityLat, cityLon)
	if err != nil {
		s.Logger.Error("Failed to insert/select city: ", err)
		return nil, err
	}
	defer row.Close()

	city, err := pgx.CollectOneRow(row, pgx.RowToStructByName[City])
	if err != nil {
		s.Logger.Error("Failed to collect city: ", err)
		return nil, err
	}

	buildingQuery := `INSERT INTO geo.buildings (city_id, street, additional, house_number, lat, lon)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, city_id, street, additional, house_number, lat, lon`
	row, err = tx.Query(ctx, buildingQuery, city.ID, street, additional, houseNumber, streetLat, streetLon)
	if err != nil {
		s.Logger.Error("Failed to insert/select building: ", err)
		return nil, err
	}
	defer row.Close()

	building, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Building])
	if err != nil {
		s.Logger.Error("Failed to collect building: ", err)
		return nil, err
	}

	addressQuery := `INSERT INTO geo.addresses (building_id, apartment_number)
	VALUES ($1, $2)
	ON CONFLICT (building_id, apartment_number) DO NOTHING
	RETURNING id, building_id, apartment_number`
	row, err = tx.Query(ctx, addressQuery, building.ID, apartmentNumber)
	if err != nil {
		s.Logger.Error("Failed to insert/select address: ", err)
		return nil, err
	}
	defer row.Close()

	address, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Address])
	if err != nil {
		s.Logger.Error("Failed to collect address: ", err)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		s.Logger.Error("Failed to commit transaction: ", err)
		return nil, err
	}

	fullAddress := city.Name + ", " + building.Street + ", " + building.HouseNumber
	if address.ApartmentNumber != "" {
		fullAddress += ", " + address.ApartmentNumber
	}

	return &FullAddress{
		RegionID:        region.ID,
		Region:          region.Name,
		RegionLat:       region.Lat,
		RegionLon:       region.Lon,
		CityID:          city.ID,
		City:            city.Name,
		CityLat:         city.Lat,
		CityLon:         city.Lon,
		BuildingID:      building.ID,
		Street:          building.Street,
		StreetLat:       building.Lat,
		StreetLon:       building.Lon,
		HouseNumber:     building.HouseNumber,
		AddressID:       address.ID,
		ApartmentNumber: address.ApartmentNumber,
		FullAddress:     fullAddress,
	}, nil
}

func NewGeoCRUD(db *database.DataBase, log *logrus.Logger) *CRUDGeo {
	return &CRUDGeo{
		DataBase: db,
		Logger:   log,
	}
}
