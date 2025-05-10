package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Census-Population-Project/API/internal/service/api/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(v)
	if err != nil {
		var jsonErr *json.UnmarshalTypeError
		if errors.As(err, &jsonErr) {
			return NewInvalidValueForFieldError(jsonErr.Field)
		}

		if strings.Contains(err.Error(), "json: unknown field") {
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return NewUnknownFieldError(strings.Replace(field, `"`, ``, -1))
		}

		return err
	}
	return nil
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, response.ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func ParseStringToList(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, ",")
}

func ParseIntQuery(r *http.Request, key string, min, max int, defaultValue int, nilDefaultValue bool) (*int, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		if nilDefaultValue {
			return nil, nil
		}
		return &defaultValue, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil || val < min || val > max {
		return nil, fmt.Errorf("Invalid %s", key)
	}
	return &val, nil
}

func ParseInt64Query(r *http.Request, key string, min, max int64, defaultValue int64, nilDefaultValue bool) (*int64, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		if nilDefaultValue {
			return nil, nil
		}
		return &defaultValue, nil
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil || val < min || val > max {
		return nil, fmt.Errorf("Invalid %s", key)
	}
	return &val, nil
}

func ParseStringQuery(r *http.Request, key string) (*string, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		//return nil, fmt.Errorf("Invalid %s", key)
		return nil, nil
	}
	return &valStr, nil
}

func StringTimeToTimeWithTZ(timeString string, format TimeFormat) (*time.Time, error) {
	parsedTime, err := time.Parse(string(format), timeString)
	if err != nil {
		return nil, NewInvalidTimeFormatError()
	}
	parsedTime = parsedTime.UTC()
	return &parsedTime, nil
}

func GetUserIdFromContext(r *http.Request) (*uuid.UUID, error) {
	claims, ok := r.Context().Value("claims").(*jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	userID, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id not found in claims")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format")
	}
	return &uid, nil
}

func UpdateOptionalField[T any](opt Optional[T], current *T, isPatch bool, nonNil bool) *T {
	if opt.Value != nil {
		return opt.Value
	} else {
		if isPatch {
			if opt.Defined {
				if nonNil {
					return current
				}
				return nil
			} else {
				return current
			}
		} else {
			if nonNil {
				return current
			}
			return nil
		}
	}
}
