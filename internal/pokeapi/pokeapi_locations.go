package pokeapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *PokeAPIClient) getLocationAreas(URL string) (ShallowLocations, error) {
	if URL == "" {
		URL = c.baseUrl + "location-area/"
	}
	locationState := &c.paginationStates.LocationState

	body, found := c.cache.Get(URL)
	if !found {
		resp, err := http.Get(URL)
		if err != nil {
			return ShallowLocations{}, fmt.Errorf("failed to connect to PokeAPI: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return ShallowLocations{}, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return ShallowLocations{}, fmt.Errorf("Failed to read response body: %w", err)
		}
		c.cache.Add(URL, body)
	}

	var response ShallowLocations
	err := json.Unmarshal(body, &response)
	if err != nil {
		return ShallowLocations{}, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}

	locationState.NextURL = response.NextURL
	locationState.PreviousURL = response.PreviousURL

	return response, nil
}

func (c *PokeAPIClient) GetNextLocationAreas() (ShallowLocations, error) {
	locationState := &c.paginationStates.LocationState
	if locationState.NextURL == nil && locationState.PreviousURL == nil {
		return c.getLocationAreas("")
	}
	URL, err := locationState.GoForward()
	if err != nil {
		return ShallowLocations{}, err
	}
	return c.getLocationAreas(URL)
}

func (c *PokeAPIClient) GetPrevLocationAreas() (ShallowLocations, error) {
	locationState := &c.paginationStates.LocationState
	if locationState.NextURL == nil && locationState.PreviousURL == nil {
		return c.getLocationAreas("")
	}
	URL, err := locationState.GoBack()
	if err != nil {
		return ShallowLocations{}, err
	}
	return c.getLocationAreas(URL)
}

// Get pokemon in location

func (c *PokeAPIClient) GetNamedLocation(name string) (NamedLocationArea, error) {
	URL := c.baseUrl + "location-area/" + name

	body, found := c.cache.Get(URL)
	if !found {
		resp, err := http.Get(URL)
		if err != nil {
			return NamedLocationArea{}, fmt.Errorf("failed to connect to PokeAPI: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return NamedLocationArea{}, fmt.Errorf("That is not a valid region, use map to find one.")
		}
		if resp.StatusCode != 200 {
			return NamedLocationArea{}, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return NamedLocationArea{}, fmt.Errorf("Failed to read response body: %w", err)
		}
		c.cache.Add(URL, body)
	}

	var response NamedLocationArea
	err := json.Unmarshal(body, &response)
	if err != nil {
		return NamedLocationArea{}, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}

	return response, nil
}
