package pokeapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/StupidWeasel/bootdev-pokedex/internal/pokecache"
	"time"
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

type PokeAPIClient struct {
    cache *pokecache.Cache
    baseUrl string
    paginationStates PaginationStates
}
func NewPokeAPIClient(interval time.Duration) *PokeAPIClient {
    return &PokeAPIClient{
        cache: pokecache.NewCache(interval * time.Minute),
        baseUrl: "https://pokeapi.co/api/v2/location-area",
        paginationStates: PaginationStates{
            LocationState: Pagination{
                NextURL: nil,
                PreviousURL: nil,
            },
        },
    }
}

func (c *PokeAPIClient) getLocationAreas(URL string) ([]LocationArea, error){
	if URL == ""{
		URL = c.baseUrl
	}
	locationState := &c.paginationStates.LocationState

	body, found := c.cache.Get(URL);
	if !found {
		//fmt.Println("[Uncached]")
		resp, err := http.Get(URL)
        if err != nil {
            return nil, fmt.Errorf("failed to connect to PokeAPI: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            return nil, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
        }

        body, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, fmt.Errorf("Failed to read response body: %w", err)
        }
        c.cache.Add(URL, body)
	}else{
		//fmt.Println("[Cached]")
	}

	var response PaginatedResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}

	locationState.NextURL = response.NextURL
	locationState.PreviousURL = response.PreviousURL

	return response.Results, nil
}

func (c *PokeAPIClient) GetNextLocationAreas() ([]LocationArea, error){
	locationState := &c.paginationStates.LocationState
	if locationState.NextURL == nil && locationState.PreviousURL == nil{
		return c.getLocationAreas("")
	} 
	URL,err := locationState.GoForward()
	if err != nil {
		return nil, err
	}
	return c.getLocationAreas(URL)
}

func (c *PokeAPIClient) GetPrevLocationAreas() ([]LocationArea, error){
	locationState := &c.paginationStates.LocationState
	if locationState.NextURL == nil && locationState.PreviousURL == nil{
		return c.getLocationAreas("")
	} 
	URL,err := locationState.GoBack()
	if err != nil {
		return nil, err
	}
	return c.getLocationAreas(URL)
}