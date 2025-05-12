package geo

import (
	"net/http"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"
)

// RegionNotFoundError is an error that is returned when a region is not found.
type regionNotFoundError struct{}

func (e *regionNotFoundError) ErrorStatusCode() int { return http.StatusNotFound }

func (*regionNotFoundError) Error() string { return "Region not found" }

func NewRegionNotFoundError() serviceerrors.ServiceError {
	return &regionNotFoundError{}
}

// RegionAlreadyExistsError is an error that is returned when a region already exists.
type regionAlreadyExistsError struct{}

func (e *regionAlreadyExistsError) ErrorStatusCode() int { return http.StatusConflict }

func (*regionAlreadyExistsError) Error() string { return "Region already exists" }

func NewRegionAlreadyExistsError() serviceerrors.ServiceError {
	return &regionAlreadyExistsError{}
}

// CityNotFoundError is an error that is returned when a city is not found.
type cityNotFoundError struct{}

func (e *cityNotFoundError) ErrorStatusCode() int { return http.StatusNotFound }

func (*cityNotFoundError) Error() string { return "City not found" }

func NewCityNotFoundError() serviceerrors.ServiceError {
	return &cityNotFoundError{}
}

// CityAlreadyExistsError is an error that is returned when a city already exists.
type cityAlreadyExistsError struct{}

func (e *cityAlreadyExistsError) ErrorStatusCode() int { return http.StatusConflict }

func (*cityAlreadyExistsError) Error() string { return "City already exists" }

func NewCityAlreadyExistsError() serviceerrors.ServiceError {
	return &cityAlreadyExistsError{}
}

// BuildingNotFoundError is an error that is returned when a building is not found.
type buildingNotFoundError struct{}

func (e *buildingNotFoundError) ErrorStatusCode() int { return http.StatusNotFound }

func (*buildingNotFoundError) Error() string { return "Building not found" }

func NewBuildingNotFoundError() serviceerrors.ServiceError {
	return &buildingNotFoundError{}
}

// BuildingAlreadyExistsError is an error that is returned when a building already exists.
type buildingAlreadyExistsError struct{}

func (e *buildingAlreadyExistsError) ErrorStatusCode() int { return http.StatusConflict }

func (*buildingAlreadyExistsError) Error() string { return "Building already exists" }

func NewBuildingAlreadyExistsError() serviceerrors.ServiceError {
	return &buildingAlreadyExistsError{}
}
