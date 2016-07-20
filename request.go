package kkcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

var (
	// ErrNoPlace to define a place can't be detected, such as in the ocean.
	ErrNoPlace = errors.New("No place found.")

	// ErrLimitation to define request over the limitation.
	ErrLimitation = errors.New("Request too many.")
)

// googleKey used to request the APIs.
var googleKey string

type statusField struct {
	Status string `json:"status"`
}

type oneAddress struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type oneLatLngResult struct {
	AddressComponents json.RawMessage `json:"address_components"`
	Formatted         json.RawMessage `json:"formatted_address"`
	Geometry          json.RawMessage `json:"geometry"`
	PlaceID           string          `json:"place_id"`
}

// latLngLocation to define the result of search result with lat and lng.
type latLngLocation struct {
	Results []oneLatLngResult `json:"results"`
	statusField
}

type predictResult struct {
	ID          json.RawMessage `json:"id"`
	Reference   json.RawMessage `json:"reference"`
	Matched     json.RawMessage `json:"matched_substrings"`
	Description string          `json:"description"`
	PlaceID     string          `json:"place_id"`
	Types       json.RawMessage `json:"types"`
	Terms       json.RawMessage `json:"terms"`
}

type predictLocation struct {
	Results []predictResult `json:"predictions"`
	statusField
}

type placeDetailResult struct {
	AddressComponents []oneAddress `json:"address_components"`
	Address           string       `json:"formatted_address"`
}

type placeDetailResponse struct {
	HTML    json.RawMessage   `json:"html_attributions"`
	Results placeDetailResult `json:"result"`
	statusField
}

// return the name and whether existed.
func getString(addr []oneAddress, tp string, isShort bool) (string, bool) {
	for _, one := range addr {
		for _, str := range one.Types {
			if str == tp {
				if isShort {
					return one.ShortName, true
				}
				return one.LongName, true
			}
		}
	}
	return "", false
}

// getLocationWithLatLng to get location with lat lng.
// If no result, return ErrNoPlace.
// If out of limitation, return ErrLimitation
func requestLocationWithLatLng(lat, lng float32) (string, error) {
	request := gorequest.New().Timeout(10 * time.Second)
	request.Type("json")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?result_type=locality&key=%s&latlng=%f,%f", googleKey, lat, lng)

	if resp, body, errs := request.Get(url).EndBytes(); len(errs) != 0 {
		return "", errs[0]
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf("Response status: %d", resp.StatusCode)
	} else {
		var result latLngLocation
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}

		if result.Status == "ZERO_RESULTS" {
			return "", ErrNoPlace
		} else if result.Status == "OVER_QUERY_LIMIT" {
			return "", ErrLimitation
		} else if result.Status == "OK" {
			return result.Results[0].PlaceID, nil
		}

		return "", errors.New("Unhandled result.")
	}
}

// getAutoComplete to get placeids and their description with input.
// Return place ids, descriptions, error
func requestAutoComplete(input, lang string) ([]string, []string, error) {
	request := gorequest.New().Timeout(10 * time.Second)
	request.Type("json")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/autocomplete/json?types=(cities)&language=%s&key=%s&input=%s", lang, googleKey, input)

	if resp, body, errs := request.Get(url).EndBytes(); len(errs) != 0 {
		return nil, nil, errs[0]
	} else if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("Response status: %d", resp.StatusCode)
	} else {
		var result predictLocation
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, nil, err
		}

		if result.Status == "ZERO_RESULTS" {
			return nil, nil, ErrNoPlace
		} else if result.Status == "OVER_QUERY_LIMIT" {
			return nil, nil, ErrLimitation
		} else if result.Status == "OK" {
			var placeids []string
			var descriptions []string

			for _, one := range result.Results {
				placeids = append(placeids, one.PlaceID)
				descriptions = append(descriptions, one.Description)
			}

			return placeids, descriptions, nil
		}

		return nil, nil, errors.New("Unhandled result.")
	}
}

// getPlaceInfo to get place information with place ID.
func requestPlaceInfo(placeid, lang string) (country, countryName, placeName, address string, erro error) {
	request := gorequest.New().Timeout(10 * time.Second)
	request.Type("json")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?placeid=%s&key=%s&language=%s", placeid, googleKey, lang)

	if resp, body, errs := request.Get(url).EndBytes(); len(errs) != 0 {
		erro = errs[0]
	} else if resp.StatusCode != 200 {
		erro = fmt.Errorf("Response status: %d", resp.StatusCode)
	} else {
		var result placeDetailResponse
		if err := json.Unmarshal(body, &result); err != nil {
			erro = err
			return
		}

		if result.Status == "ZERO_RESULTS" {
			erro = ErrNoPlace
		} else if result.Status == "OVER_QUERY_LIMIT" {
			erro = errors.New("Request too many.")
		} else if result.Status == "OK" {
			country, _ = getString(result.Results.AddressComponents, "country", true)
			countryName, _ = getString(result.Results.AddressComponents, "country", false)
			placeName, _ = getString(result.Results.AddressComponents, "locality", true)
			address = result.Results.Address
		}
	}
	return
}
