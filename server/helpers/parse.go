package helpers

import (
	"fmt"
	"net/http"
	"strconv"
)

func RequiredParamInt(r *http.Request, w http.ResponseWriter, param string) (int, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		http.Error(w, "Missing required parameter: "+param, http.StatusBadRequest)
		return -1, false
	}

	intVal, err := strconv.Atoi(values[0])
	if err != nil {
		http.Error(w, "Invalid value for "+param, http.StatusBadRequest)
		return -1, false
	}
	return intVal, true
}
func OptionalParamInt(r *http.Request, param string) (*int, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		return nil, false // Parameter not present
	}

	intVal, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, false // Invalid integer value
	}

	return &intVal, true
}

func RequiredParamString(r *http.Request, w http.ResponseWriter, param string) (string, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		http.Error(w, "Missing required parameter: "+param, http.StatusBadRequest)
		return "", false
	}

	return values[0], true
}
func OptionalParamString(r *http.Request, param string) (string, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		return "", false // Parameter not present
	}

	return values[0], true
}

func RequiredParamBool(r *http.Request, w http.ResponseWriter, param string) (bool, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		http.Error(w, fmt.Sprintf("Missing required parameter: %s", param), http.StatusBadRequest)
		return false, false
	}

	boolVal, err := strconv.ParseBool(values[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid boolean parameter: %s", param), http.StatusBadRequest)
		return false, false
	}

	return boolVal, true
}

func OptionalParamBool(r *http.Request, param string) (bool, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		return false, false // Parameter not provided
	}

	boolVal, err := strconv.ParseBool(values[0])
	if err != nil {
		return false, false // Invalid boolean format
	}

	return boolVal, true
}

func RequiredParamFloat(r *http.Request, w http.ResponseWriter, param string) (float64, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		http.Error(w, fmt.Sprintf("Missing required parameter: %s", param), http.StatusBadRequest)
		return 0, false
	}

	floatVal, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid float parameter: %s", param), http.StatusBadRequest)
		return 0, false
	}

	return floatVal, true
}

func OptionalParamFloat(r *http.Request, param string) (float64, bool) {
	values := r.URL.Query()[param]
	if len(values) == 0 {
		return 0, false // Parameter not provided
	}

	floatVal, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return 0, false // Invalid float format
	}

	return floatVal, true
}
