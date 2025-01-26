package pokeapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PaginatedResponse struct {
	Results    []LocationArea `json:"results"`
	NextURL    *string        `json:"next"`
	PreviousURL *string       `json:"previous"`
}

type Pagination struct {
	NextURL     *string
	PreviousURL *string
}

type PaginationStates struct {
	LocationState Pagination
}

var paginationStates = PaginationStates{
	LocationState: Pagination{
		NextURL:     nil,
		PreviousURL: nil,
	},
}

func (ps *PaginationStates) ResetLocationPagination(){
	ps.LocationState = Pagination{
		NextURL:     nil,
		PreviousURL: nil,
	}
}

func (p *Pagination) GoForward() (string, error){
	if p.NextURL == nil {
		return "", fmt.Errorf("You're on the last page!")
	}
	return *p.NextURL, nil
}

func (p *Pagination) GoBack() (string, error){
	if p.PreviousURL == nil {
		return "", fmt.Errorf("You're on the first page!")
	}
	return *p.PreviousURL, nil
}

func getLocationAreas(URL string) ([]LocationArea, error){
	if URL == ""{
		URL = "https://pokeapi.co/api/v2/location-area"
	}
	locationState := &paginationStates.LocationState

	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PokeAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	var response PaginatedResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}

	locationState.NextURL = response.NextURL
	locationState.PreviousURL = response.PreviousURL

	return response.Results, nil
}

func GetNextLocationAreas() ([]LocationArea, error){
	locationState := &paginationStates.LocationState
	if locationState.NextURL == nil && locationState.PreviousURL == nil{
		return getLocationAreas("")
	} 
	URL,err := locationState.GoForward()
	if err != nil {
		return nil, err
	}
	return getLocationAreas(URL)
}

func GetPrevLocationAreas() ([]LocationArea, error){
	locationState := &paginationStates.LocationState
	if locationState.NextURL == nil && locationState.PreviousURL == nil{
		return getLocationAreas("")
	} 
	URL,err := locationState.GoBack()
	if err != nil {
		return nil, err
	}
	return getLocationAreas(URL)
}