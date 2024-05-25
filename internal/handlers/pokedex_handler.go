package handlers

import (
	"encoding/json"
	"github.com/mdhuy17/project_netcentric_g5/internal/repositories"
	"github.com/mdhuy17/project_netcentric_g5/utils"
)

type PokeDexHandler struct {
	PokedexReposiotry *repositories.PokedexRepository
}

func NewPokeDexHandler(pokedexReposiotry *repositories.PokedexRepository) *PokeDexHandler {
	return &PokeDexHandler{
		PokedexReposiotry: pokedexReposiotry,
	}
}

func (s *PokeDexHandler) GetPokemon(name string) ([]byte, error) {
	data, err := s.PokedexReposiotry.GetMonsterByID(utils.PokeMap[name])
	if err != nil {
		return []byte{}, err
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return []byte{}, err
	}
	return jsonData, nil
}
