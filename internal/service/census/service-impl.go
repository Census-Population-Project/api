package census

import (
	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Config *config.Config

	DB     *database.DataBase
	Logger *logrus.Logger

	DDsAPI *suggest.Api
}

func NewService(cfg *config.Config, db *database.DataBase, logger *logrus.Logger, ddsApi *suggest.Api) *Service {
	return &Service{
		Config: cfg,
		DB:     db,
		Logger: logger,

		DDsAPI: ddsApi,
	}
}
