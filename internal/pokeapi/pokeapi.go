package pokeapi

import (
	"github.com/StupidWeasel/bootdev-pokedex/internal/pokecache"
	"time"
)

func NewPokeAPIClient(interval time.Duration) *PokeAPIClient {
    return &PokeAPIClient{
        cache: pokecache.NewCache(interval),
        baseUrl: "https://pokeapi.co/api/v2/",
        paginationStates: PaginationStates{
            LocationState: Pagination{
                NextURL: nil,
                PreviousURL: nil,
            },
        },
        Pokedex: make(map[string]Pokemon,0),
    }
}