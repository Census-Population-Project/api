package geo

import (
	"github.com/Census-Population-Project/API/internal/database"

	"github.com/sirupsen/logrus"
)

type CRUDInterface interface {
}

type CRUDGeo struct {
	DataBase *database.DataBase
	Logger   *logrus.Logger
}

func NewGeoCRUD(db *database.DataBase, log *logrus.Logger) *CRUDGeo {
	return &CRUDGeo{
		DataBase: db,
		Logger:   log,
	}
}
