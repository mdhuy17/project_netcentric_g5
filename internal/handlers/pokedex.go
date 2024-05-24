package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mdhuy17/project_netcentric_g5/internal/models/api/v1/description"
	_type "github.com/mdhuy17/project_netcentric_g5/internal/models/api/v1/type"
	"github.com/mdhuy17/project_netcentric_g5/internal/models/evolutions"
	"github.com/mdhuy17/project_netcentric_g5/internal/models/monster_moves"
	"github.com/mdhuy17/project_netcentric_g5/internal/models/monster_supplementals"
	"github.com/mdhuy17/project_netcentric_g5/internal/models/moves"
	"github.com/mdhuy17/project_netcentric_g5/internal/models/skim_monsters"
	"github.com/mdhuy17/project_netcentric_g5/utils"
	"io/ioutil"
)

type Pokedex struct {
}

var BasePath = "./internal/models"

type Pokemon struct {
	Skim                *skim_monsters.Data
	Evol                *evolutions.Data
	MonsterMoves        []*moves.Data
	MonsterSupplemental *monster_supplementals.Data
	Description         []*description.Data
	Types               []*_type.Data
}

func (p *Pokedex) GetMonsterMovesByID(id string) ([]*moves.Data, error) {
	var data []*moves.Data
	pathFile := fmt.Sprintf("%s/monster_moves/data/%s.json", BasePath, id)
	file, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var monsterMove monster_moves.Data
	err = json.Unmarshal(file, &monsterMove)
	if err != nil {
		return nil, err
	}

	requestMoves := monsterMove.Move
	for _, move := range requestMoves {
		pathFile = fmt.Sprintf("%s/moves/data/%d.json", BasePath, move.Id)
		file, err = ioutil.ReadFile(pathFile)
		if err != nil {
			return nil, err
		}
		var m moves.Data
		err = json.Unmarshal(file, &m)
		if err != nil {
			return nil, err
		}
		data = append(data, &m)
	}
	return data, nil

}

func (p *Pokedex) GetMonsterTypeByID(path []skim_monsters.ListMapObject) ([]*_type.Data, error) {
	var data []*_type.Data
	for _, id := range path {
		pathFile := fmt.Sprintf("%s/api/v1/type/%s/poke.json", BasePath, id.Name)
		file, err := ioutil.ReadFile(pathFile)
		if err != nil {
			return nil, err

		}
		var t _type.Data
		err = json.Unmarshal(file, &t)
		if err != nil {
			return nil, err

		}
		data = append(data, &t)
	}
	return data, nil
}

func (p *Pokedex) GetMonsterDescription(path []skim_monsters.ListMapObject) ([]*description.Data, error) {
	var data []*description.Data
	for _, id := range path {
		pathFile := fmt.Sprintf("%s%s/poke.json", BasePath, id.ResourceURI)
		file, err := ioutil.ReadFile(pathFile)
		if err != nil {
			return nil, err
		}
		var desc description.Data
		err = json.Unmarshal(file, &desc)
		if err != nil {
			return nil, err
		}
		data = append(data, &desc)
	}

	return data, nil

}

func (p *Pokedex) GetMonsterByID(id string) (*Pokemon, error) {
	pathFile := fmt.Sprintf("%s/skim_monsters/data/%s.json", BasePath, id)
	file, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var monster skim_monsters.Data
	err = json.Unmarshal(file, &monster)
	if err != nil {
		return nil, err
	}

	pathFile = fmt.Sprintf("%s/evolutions/data/%s.json", BasePath, id)
	file, err = ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var evol evolutions.Data
	err = json.Unmarshal(file, &evol)
	if err != nil {
		return nil, err
	}

	pathFile = fmt.Sprintf("%s/monster_supplementals/data/%s.json", BasePath, id)
	file, err = ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var supp monster_supplementals.Data
	err = json.Unmarshal(file, &supp)
	if err != nil {
		return nil, err
	}

	monsterMoves, err := p.GetMonsterMovesByID(id)
	if err != nil {
		return nil, err
	}

	desc, err := p.GetMonsterDescription(monster.Description)
	if err != nil {
		return nil, err

	}

	types, err := p.GetMonsterTypeByID(monster.Type)
	if err != nil {
		return nil, err
	}

	return &Pokemon{
		Skim:                &monster,
		Evol:                &evol,
		MonsterMoves:        monsterMoves,
		MonsterSupplemental: &supp,
		Description:         desc,
		Types:               types,
	}, nil
}

func main() {

	var pokedex Pokedex

	data, err := pokedex.GetMonsterByID(utils.PokeMap["Ivysaur"])
	if err != nil {
		return
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print JSON data
	fmt.Println(string(jsonData))

}
