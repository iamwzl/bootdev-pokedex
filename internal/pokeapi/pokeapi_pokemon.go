package pokeapi

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *PokeAPIClient) GetNamedPokemon(name string) (Pokemon, error){
	


	URL := c.baseUrl + "pokemon/" + name

	body, found := c.cache.Get(URL);
	if !found {
		resp, err := http.Get(URL)
        if err != nil {
            return Pokemon{}, fmt.Errorf("failed to connect to PokeAPI: %w", err)
        }
        defer resp.Body.Close()
        if resp.StatusCode == 404{
        	return Pokemon{}, fmt.Errorf("That is not a valid pokemon.")
        }
        if resp.StatusCode != 200 {
            return Pokemon{}, fmt.Errorf("PokeAPI returned status code %d", resp.StatusCode)
        }

        body, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            return Pokemon{}, fmt.Errorf("Failed to read response body: %w", err)
        }
        c.cache.Add(URL, body)
	}

	var response Pokemon
	err := json.Unmarshal(body, &response)
	if err != nil {
		return Pokemon{}, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}

	return response, nil
}
