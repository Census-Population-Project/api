package geo

import (
	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Config *config.Config

	DB     *database.DataBase
	Logger *logrus.Logger

	CRUDGeo *CRUDGeo
}

func NewService(cfg *config.Config, db *database.DataBase, logger *logrus.Logger) *Service {
	return &Service{
		Config: cfg,

		DB:     db,
		Logger: logger,

		CRUDGeo: NewGeoCRUD(db, logger),
	}
}
