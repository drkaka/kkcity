package request

import "encoding/json"

type statusField struct {
	Status string `json:"status"`
}

type oneAddress struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type oneLatlngResult struct {
	AddressComponents []oneAddress    `json:"address_components"`
	Formatted         string          `json:"formatted_address"`
	Geometry          json.RawMessage `json:"geometry"`
	PlaceID           string          `json:"place_id"`
}

// LatLngLocation to define the result of search result with lat and lng.
type latLngLocation struct {
	Results []oneLatlngResult `json:"results"`
	statusField
}

type predictTerm struct {
	Offset int32  `json:"offset"`
	Value  string `json:"value"`
}

type predictResult struct {
	ID          json.RawMessage `json:"id"`
	Reference   json.RawMessage `json:"reference"`
	Matched     json.RawMessage `json:"matched_substrings"`
	Description string          `json:"description"`
	PlaceID     string          `json:"place_id"`
	Types       json.RawMessage `json:"types"`
	Terms       []predictTerm   `json:"terms"`
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
