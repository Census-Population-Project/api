package census

import (
	"github.com/Census-Population-Project/API/internal/database"

	"github.com/sirupsen/logrus"
)

type CRUDInterface interface{}

type CRUDCensus struct {
	DataBase *database.DataBase
	Logger   *logrus.Logger
}

func NewCensusCRUD(db *database.DataBase, log *logrus.Logger) *CRUDCensus {
	return &CRUDCensus{
		DataBase: db,
		Logger:   log,
	}
}
