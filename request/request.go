package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

// GoogleKey used to request the APIs.
var GoogleKey string

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

// GetLocationWithLatLng to get location with lat lng.
// If no result, return error ErrNoPlace.
func GetLocationWithLatLng(lat, lng float32, language string) (*LatLngResult, error) {
	request := gorequest.New().Timeout(10 * time.Second)
	request.Type("json")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?result_type=locality&key=%s&latlng=%f,%f&language=%s", GoogleKey, lat, lng, language)

	//"OK" means there will be at least one result.
	//"ZERO_RESULTS" successed to geoencoding but no result.
	//"OVER_QUERY_LIMIT" out of limitation.
	if resp, body, errs := request.Get(url).EndBytes(); len(errs) != 0 {
		return nil, errs[0]
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Response status: %d", resp.StatusCode)
	} else {
		var result latLngLocation
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		if result.Status == "ZERO_RESULTS" {
			return nil, ErrNoPlace
		} else if result.Status == "OVER_QUERY_LIMIT" {
			return nil, errors.New("Request too many.")
		} else if result.Status == "OK" {
			var realResult LatLngResult

			one := result.Results[0]
			realResult.Description = one.Formatted
			realResult.PlaceID = one.PlaceID
			realResult.City, _ = getString(one.AddressComponents, "locality", true)
			return &realResult, nil
		}

		return nil, errors.New("Unhandled result.")
	}
}
