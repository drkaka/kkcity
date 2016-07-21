package kkcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLocationWithLatLng(t *testing.T) {
	placeid, err := requestLocationWithLatLng(24.54918, 118.12705)
	assert.NoError(t, err, "Should get the city information.")
	assert.Equal(t, "ChIJJ-u_5XmDFDQRVtBolgpnoCg", placeid, "Place ID result is wrong.")

	_, err = requestLocationWithLatLng(0, 0)
	assert.Equal(t, ErrNoPlace, err, "Should find no place.")
}

func TestRequestAutoComplete(t *testing.T) {
	placeids, descriptions, err := requestAutoComplete("bao", "en")
	assert.Nil(t, err, "Shoule be able to get auto complete result.")

	assert.EqualValues(t, 5, len(placeids), "Should get max record.")
	assert.EqualValues(t, 5, len(descriptions), "Should get max record.")

	// 0
	assert.Equal(t, "Baotou, Inner Mongolia, China", descriptions[0], "Description result is wrong.")
	assert.Equal(t, "ChIJ2xkwFyJYBDYRfF3XPzwYEZo", placeids[0], "Place ID result is wrong.")

	// 1
	assert.Equal(t, "Baoding, Hebei, China", descriptions[1], "Description result is wrong.")
	assert.Equal(t, "ChIJlQRfqI2P5TURK0rev-lCO84", placeids[1], "Place ID result is wrong.")

	// 2
	assert.Equal(t, "Baoji, Shaanxi, China", descriptions[2], "Description result is wrong.")
	assert.Equal(t, "ChIJ6ZCbaAPdYDYRVyMhVjeyeOs", placeids[2], "Place ID result is wrong.")

	// 3
	assert.Equal(t, "Baoshan, Yunnan, China", descriptions[3], "Description result is wrong.")
	assert.Equal(t, "ChIJkc6IxzA2LzcReDKMDHoyouE", placeids[3], "Place ID result is wrong.")
}

func TestRequestPlaceInfo(t *testing.T) {
	placeid := "ChIJJ-u_5XmDFDQRVtBolgpnoCg"
	country, countryName, placeName, address, err := requestPlaceInfo(placeid, "en")
	assert.Nil(t, err, "Should be able to get place information.")
	assert.Equal(t, "CN", country, "Country information wrong.")
	assert.Equal(t, "China", countryName, "Country name information wrong.")
	assert.Equal(t, "Xiamen", placeName, "Place name information wrong.")
	assert.Equal(t, "Xiamen, Fujian, China", address, "Address information wrong.")

	country, countryName, placeName, address, err = requestPlaceInfo(placeid, "zh")
	assert.Nil(t, err, "Should be able to get place information.")
	assert.Equal(t, "CN", country, "Country information wrong.")
	assert.Equal(t, "中国", countryName, "Country name information wrong.")
	assert.Equal(t, "厦门", placeName, "Place name information wrong.")
	assert.Equal(t, "中国福建省厦门市", address, "Address information wrong.")
}
