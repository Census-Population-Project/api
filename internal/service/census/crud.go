package census

import (
	"context"

	"github.com/Census-Population-Project/API/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type CRUDInterface interface {
	SelectEvents(limit, offset int) ([]Event, *int64, error)
	SelectEventInfoByID(id uuid.UUID) (*EventInfo, error)
	SelectEventInfoInLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventDataInLocation, error)
	SelectEventStatisticsInLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventStatistics, error)
}

type CRUDCensus struct {
	DataBase *database.DataBase
	Logger   *logrus.Logger
}

func (c *CRUDCensus) SelectEvents(limit, offset int) ([]Event, *int64, error) {
	query := `SELECT id, name, start_datetime, end_datetime FROM census.events ORDER BY start_datetime ASC LIMIT $1 OFFSET $2`
	rows, err := c.DataBase.DBPool.Query(context.Background(), query, limit, offset)
	if err != nil {
		c.Logger.Error("Error selecting events: ", err)
		return nil, nil, err
	}
	defer rows.Close()

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[Event])
	if err != nil {
		c.Logger.Error("Error collecting events: ", err)
		return nil, nil, err
	}

	var total *int64
	query = `SELECT COUNT(*) FROM census.events`
	row := c.DataBase.DBPool.QueryRow(context.Background(), query)
	if err := row.Scan(&total); err != nil {
		c.Logger.Error("Error counting events: ", err)
		return nil, nil, err
	}

	return events, total, nil
}

func (c *CRUDCensus) SelectEventInfoByID(id uuid.UUID) (*EventInfo, error) {
	query := `WITH persons_filtered AS (SELECT p.*, h.address_id, h.event_id, b.city_id, c.region_id
							  FROM census.persons p
									   JOIN census.households h ON h.id = p.household_id
									   JOIN geo.addresses a ON a.id = h.address_id
									   JOIN geo.buildings b ON b.id = a.building_id
									   JOIN geo.cities c ON c.id = b.city_id
							  WHERE h.event_id = $1),
		 genders_agg AS (SELECT p.gender::text AS type, COUNT(*) AS count
						 FROM census.persons p
								  JOIN census.households h ON h.id = p.household_id
						 WHERE h.event_id = $1
						 GROUP BY p.gender),
		 regions_agg AS (SELECT DISTINCT c.region_id AS id, rg.name, rg.lat, rg.lon
						 FROM census.persons p
								  JOIN census.households h ON h.id = p.household_id
								  JOIN geo.addresses a ON a.id = h.address_id
								  JOIN geo.buildings b ON b.id = a.building_id
								  JOIN geo.cities c ON c.id = b.city_id
								  JOIN geo.regions rg ON rg.id = c.region_id
						 WHERE h.event_id = $1),
		 cities_agg AS (SELECT DISTINCT b.city_id AS id, ct.region_id, ct.name, ct.lat, ct.lon
						FROM census.persons p
								 JOIN census.households h ON h.id = p.household_id
								 JOIN geo.addresses a ON a.id = h.address_id
								 JOIN geo.buildings b ON b.id = a.building_id
								 JOIN geo.cities ct ON ct.id = b.city_id
						WHERE h.event_id = $1)
	SELECT e.id,
		   e.name,
		   e.start_datetime,
		   e.end_datetime,
		   COUNT(DISTINCT pf.region_id)                                         AS regions_count,
		   COUNT(DISTINCT pf.city_id)                                           AS cities_count,
		   COUNT(pf.id)                                                         AS population_count,
		   COALESCE((SELECT json_agg(row_to_json(g)) FROM genders_agg g), '[]') AS genders,
		   COALESCE((SELECT json_agg(row_to_json(r)) FROM regions_agg r), '[]') AS regions,
		   COALESCE((SELECT json_agg(row_to_json(c)) FROM cities_agg c), '[]')  AS cities
	FROM census.events e
			 LEFT JOIN persons_filtered pf ON pf.event_id = e.id
	WHERE e.id = $1
	GROUP BY e.id`

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

	return &event, nil
}

func (c *CRUDCensus) SelectEventInfoInLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventDataInLocation, error) {
	query := `WITH persons_filtered AS (SELECT p.*, h.address_id, h.event_id, b.city_id, c.region_id
							  FROM census.persons p
									   JOIN census.households h ON h.id = p.household_id
									   JOIN geo.addresses a ON a.id = h.address_id
									   JOIN geo.buildings b ON b.id = a.building_id
									   JOIN geo.cities c ON c.id = b.city_id
							  WHERE h.event_id = $1
								AND ($2::uuid IS NULL OR c.region_id = $2)
								AND ($3::uuid IS NULL OR b.city_id = $3)),
		 genders_agg AS (SELECT p.gender::text AS type, COUNT(*) AS count
						 FROM census.persons p
								  JOIN census.households h ON h.id = p.household_id
								  JOIN geo.addresses a ON a.id = h.address_id
								  JOIN geo.buildings b ON b.id = a.building_id
								  JOIN geo.cities c ON c.id = b.city_id
						 WHERE h.event_id = $1
						   AND ($2::uuid IS NULL OR c.region_id = $2)
						   AND ($3::uuid IS NULL OR b.city_id = $3)
						 GROUP BY p.gender)
	SELECT e.id,
		   e.name,
		   e.start_datetime,
		   e.end_datetime,
		   COUNT(pf.id)                                                         AS population_count,
		   COALESCE((SELECT json_agg(row_to_json(g)) FROM genders_agg g), '[]') AS genders
	FROM census.events e
			 LEFT JOIN persons_filtered pf ON pf.event_id = e.id
	WHERE e.id = $1
	GROUP BY e.id`

	row, err := c.DataBase.DBPool.Query(context.Background(), query, id, regionId, cityId)
	if err != nil {
		c.Logger.Error("Failed to select event info in location IDs: ", err)
		return nil, err
	}
	defer row.Close()

	event, err := pgx.CollectOneRow(row, pgx.RowToStructByName[EventDataInLocation])
	if err != nil {
		c.Logger.Error("Failed to collect event info in location IDs: ", err)
		return nil, err
	}

	return &event, nil
}

func (c *CRUDCensus) SelectEventStatisticsInLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventStatistics, error) {
	query := `WITH filtered_addresses AS (SELECT a.id AS address_id
								FROM geo.addresses a
										 JOIN geo.buildings b ON a.building_id = b.id
										 JOIN geo.cities c ON b.city_id = c.id
										 JOIN geo.regions r ON c.region_id = r.id
								WHERE ($2::uuid IS NULL OR c.id = $2)
								  AND ($3::uuid IS NULL OR r.id = $3)),
	
		 filtered_households AS (SELECT h.*
								 FROM census.households h
								 WHERE h.event_id = $1
								   AND h.address_id IN (SELECT address_id FROM filtered_addresses)),
	
		 filtered_persons AS (SELECT p.*
							  FROM census.persons p
									   JOIN filtered_households h ON p.household_id = h.id),
	
		 age_data AS (SELECT id,
							 EXTRACT(YEAR FROM AGE(current_date, birth_date))::int AS age
					  FROM filtered_persons),
	
		 genders_agg AS (SELECT p.gender::text AS type, COUNT(*) AS count
						 FROM filtered_persons p
						 GROUP BY p.gender),
	
		 education_agg AS (SELECT education_level, COUNT(*) AS edu_count
						   FROM filtered_persons
						   GROUP BY education_level),
	
		 employment_agg AS (SELECT employment_status, COUNT(*) AS emp_count
							FROM filtered_persons
							GROUP BY employment_status)
	
	SELECT
		-- General statistics
		(SELECT COUNT(*) FROM filtered_persons)                                    AS total_population,
		(SELECT COUNT(*) FROM filtered_households)                                 AS total_households,
		COALESCE(ROUND(
						 COALESCE((SELECT COUNT(*) FROM filtered_persons), 0)::NUMERIC /
						 NULLIF((SELECT COUNT(*) FROM filtered_households), 0), 2
				 ), 0)                                                             AS avg_persons_per_household,
	
		-- Population structure
		COALESCE((SELECT json_agg(row_to_json(g)) FROM genders_agg g), '[]'::json) AS gender_distribution,
		COALESCE(ROUND(AVG(age_data.age), 1), 0)                                   AS average_age,
		COUNT(*) FILTER (WHERE age_data.age < 18)                                  AS children_count,
		COUNT(*) FILTER (WHERE age_data.age >= 65)                                 AS elderly_count,
	
		-- Education and employment
		COALESCE((SELECT jsonb_object_agg(education_level, edu_count)
				  FROM education_agg
				  WHERE education_level IS NOT NULL), '{}'::jsonb)                 AS education_distribution,
	
		COALESCE((SELECT jsonb_object_agg(employment_status, emp_count)
				  FROM employment_agg
				  WHERE employment_status IS NOT NULL), '{}'::jsonb)               AS employment_distribution,
	
		-- Language and citizenship
		COALESCE(ROUND(
						 100.0 * COUNT(*) FILTER (WHERE speaks_russian) / NULLIF(COUNT(*), 0),
						 2
				 ), 0)                                                             AS percent_speaks_russian,
	
		COUNT(*) FILTER (WHERE has_dual_citizenship = true)                        AS dual_citizenship_count,
	
		COALESCE((SELECT jsonb_agg(t)
				  FROM (SELECT lang AS type, COUNT(*) AS count
						FROM (SELECT unnest(other_languages) AS lang
							  FROM filtered_persons) AS langs
						GROUP BY lang
						ORDER BY count DESC
						LIMIT 5) t), '[]'::jsonb)                                  AS top_other_languages,
	
		-- Income sources
		COALESCE((SELECT jsonb_object_agg(source, count)
				  FROM (SELECT income_source AS source, COUNT(*) AS count
						FROM (SELECT unnest(income_sources) AS income_source
							  FROM filtered_persons) sub
						GROUP BY income_source) sub2), '{}'::jsonb)                AS income_sources_distribution
	
	FROM filtered_persons
			 JOIN age_data ON age_data.id = filtered_persons.id`

	row, err := c.DataBase.DBPool.Query(context.Background(), query, id, cityId, regionId)
	if err != nil {
		c.Logger.Error("Failed to select event statistics in location IDs: ", err)
		return nil, err
	}
	defer row.Close()

	eventStatistics, err := pgx.CollectOneRow(row, pgx.RowToStructByName[EventStatistics])
	if err != nil {
		c.Logger.Error("Failed to collect event statistics in location IDs: ", err)
		return nil, err
	}

	return &eventStatistics, nil
}

func NewCensusCRUD(db *database.DataBase, log *logrus.Logger) *CRUDCensus {
	return &CRUDCensus{
		DataBase: db,
		Logger:   log,
	}
}
