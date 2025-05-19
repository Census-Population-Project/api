package census

import (
	"context"

	"github.com/Census-Population-Project/API/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type CRUDInterface interface {
	SelectEvents(limit, offset int) ([]Event, error)
	SelectEventInfoByID(id uuid.UUID) (*EventInfo, error)
	SelectEventInfoInLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventInfo, error)
}

type CRUDCensus struct {
	DataBase *database.DataBase
	Logger   *logrus.Logger
}

func (c *CRUDCensus) SelectEvents(limit, offset int) ([]Event, error) {
	query := `SELECT id, name, start_datetime, end_datetime FROM census.events LIMIT $1 OFFSET $2`
	rows, err := c.DataBase.DBPool.Query(context.Background(), query, limit, offset)
	if err != nil {
		c.Logger.Error("Error selecting events: ", err)
		return nil, err
	}
	defer rows.Close()

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[Event])
	if err != nil {
		c.Logger.Error("Error collecting events: ", err)
		return nil, err
	}

	return events, nil
}

func (c *CRUDCensus) SelectEventInfoByID(id uuid.UUID) (*EventInfo, error) {
	query := `
		WITH persons_filtered AS (
			SELECT p.*, h.address_id, h.event_id, b.city_id, c.region_id
			FROM census.persons p
					 JOIN census.households h ON h.id = p.household_id
					 JOIN geo.addresses a ON a.id = h.address_id
					 JOIN geo.buildings b ON b.id = a.building_id
					 JOIN geo.cities c ON c.id = b.city_id
			WHERE h.event_id = $1
		)
		SELECT
			e.id,
			e.name,
			e.start_datetime,
			e.end_datetime,
			COUNT(DISTINCT pf.region_id) AS regions_count,
			COUNT(DISTINCT pf.city_id) AS cities_count,
			COUNT(pf.id) AS population,
			COALESCE(
					(
						SELECT json_agg(row_to_json(g))
						FROM (
								 SELECT p.gender::text AS type, COUNT(*) AS count
								 FROM census.persons p
										  JOIN census.households h ON h.id = p.household_id
								 WHERE h.event_id = e.id
								 GROUP BY p.gender
							 ) g
					), '[]'
			) AS genders
		FROM census.events e
				 LEFT JOIN persons_filtered pf ON pf.event_id = e.id
		WHERE e.id = $1
		GROUP BY e.id
	`

	row, err := c.DataBase.DBPool.Query(context.Background(), query, id)
	if err != nil {
		c.Logger.Error("Failed to select event info: ", err)
		return nil, err
	}
	defer row.Close()

	event, err := pgx.CollectOneRow(row, pgx.RowToStructByName[EventInfo])
	if err != nil {
		c.Logger.Error("Failed to collect event info: ", err)
		return nil, err
	}

	//genderQuery := `
	//	SELECT p.gender AS type, COUNT(*) AS count
	//	FROM census.persons p
	//	WHERE p.event_id = $1
	//	GROUP BY p.gender
	//`
	//rows, err := c.DataBase.DBPool.Query(context.Background(), genderQuery, id)
	//if err != nil {
	//	c.Logger.Error("Failed to select genders: ", err)
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//genders, err := pgx.CollectRows(rows, pgx.RowToStructByName[Gender])
	//if err != nil {
	//	c.Logger.Error("Failed to collect genders: ", err)
	//	return nil, err
	//}
	//
	//event.Genders = genders
	return &event, nil
}

func (c *CRUDCensus) SelectEventInfoInLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventInfo, error) {
	return nil, nil
}

func NewCensusCRUD(db *database.DataBase, log *logrus.Logger) *CRUDCensus {
	return &CRUDCensus{
		DataBase: db,
		Logger:   log,
	}
}
