package census

import (
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

	CRUDCensus *CRUDCensus
}

func (s *Service) GetEvents(limit, offset int) ([]Event, *int64, error) {
	return s.CRUDCensus.SelectEvents(limit, offset)
}

func (s *Service) GetEventInfoByID(id uuid.UUID) (*EventInfo, error) {
	return s.CRUDCensus.SelectEventInfoByID(id)
}

func (s *Service) GetEventInfoByLocationIDs(id uuid.UUID, regionId *uuid.UUID, cityId *uuid.UUID) (*EventDataInLocation, error) {
	return s.CRUDCensus.SelectEventInfoInLocationIDs(id, regionId, cityId)
}

func NewService(cfg *config.Config, db *database.DataBase, logger *logrus.Logger, ddsApi *suggest.Api) *Service {
	return &Service{
		Config: cfg,
		DB:     db,
		Logger: logger,

		DDsAPI: ddsApi,

		CRUDCensus: NewCensusCRUD(db, logger),
	}
}
