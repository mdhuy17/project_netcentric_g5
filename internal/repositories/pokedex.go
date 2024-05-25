package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/mdhuy17/project_netcentric_g5/internal/models"
	"github.com/mdhuy17/project_netcentric_g5/utils"
	"io/ioutil"
)

type Pokedex struct {
}

var BasePath = "./internal/models"

type Pokemon struct {
	Monster             *models.Monster             `json:"monster"`
	Description         []*models.Descriptions      `json:"description"`
	Evolution           *models.Evolution           `json:"evolution"`
	Types               []*models.Types             `json:"types"`
	MonsterMoves        []*models.Move              `json:"monster_moves"`
	MonsterSupplemental *models.MonsterSupplemental `json:"monster_supplemental"`
}

func (p *Pokedex) GetMonsterMovesByID(id string) ([]*models.Move, error) {
	var data []*models.Move
	pathFile := fmt.Sprintf("%s/monster_moves/data/%s.json", BasePath, id)
	file, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var monsterMove models.MonsterMove
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
		var m models.Move
		err = json.Unmarshal(file, &m)
		if err != nil {
			return nil, err
		}
		data = append(data, &m)
	}
	return data, nil

}

func (p *Pokedex) GetMonsterTypeByID(path []models.ListMapObject) ([]*models.Types, error) {
	var data []*models.Types
	for _, id := range path {
		pathFile := fmt.Sprintf("%s/api/v1/type/%s/poke.json", BasePath, id.Name)
		file, err := ioutil.ReadFile(pathFile)
		if err != nil {
			return nil, err

		}
		var t models.Types
		err = json.Unmarshal(file, &t)
		if err != nil {
			return nil, err

		}
		data = append(data, &t)
	}
	return data, nil
}

func (p *Pokedex) GetMonsterDescription(path []models.ListMapObject) ([]*models.Descriptions, error) {
	var data []*models.Descriptions
	for _, id := range path {
		pathFile := fmt.Sprintf("%s%s/poke.json", BasePath, id.ResourceURI)
		file, err := ioutil.ReadFile(pathFile)
		if err != nil {
			return nil, err
		}
		var desc models.Descriptions
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
	var monster models.SkimMonster
	err = json.Unmarshal(file, &monster)
	if err != nil {
		return nil, err
	}

	pathFile = fmt.Sprintf("%s/evolutions/data/%s.json", BasePath, id)
	file, err = ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var evol models.Evolution
	err = json.Unmarshal(file, &evol)
	if err != nil {
		return nil, err
	}

	pathFile = fmt.Sprintf("%s/monster_supplementals/data/%s.json", BasePath, id)
	file, err = ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var supp *models.MonsterSupplemental
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
		Monster:             monster.ToMonster(),
		Evolution:           &evol,
		MonsterMoves:        monsterMoves,
		MonsterSupplemental: supp,
		Description:         desc,
		Types:               types,
	}, nil
}

func crawl() string {

	var pokedex Pokedex

	data, err := pokedex.GetMonsterByID(utils.PokeMap["Ivysaur"])
	if err != nil {
		return ""
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Print JSON data
	return string(jsonData)

}
