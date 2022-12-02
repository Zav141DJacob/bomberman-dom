package main

import (
	"errors"
	"math/rand"
	"time"
)

type MapTile uint

const (
	MapTileTypeEmpty MapTile = iota // Empty square
	MapTileTypeWall                 // Indestructible wall
	MapTileTypeBlock                // Destructible block

	// These are also used as items
	MapTileTypeBomb
	MapTileTypeBombPowerup
	MapTileTypeFlamePowerup
	MapTileTypeSpeedPowerup
)

var ERR_INCORRECT_CORDS = errors.New("Incorrect coordinates")

// 13*11=143 map
type GameMap [MAP_WIDTH * MAP_HEIGHT]MapTile

// Checks if player can move to tile
func (gameMap *GameMap) CanMoveTo(x uint, y uint) bool {
	tile, err := gameMap.GetTile(x, y)
	return tile != MapTileTypeWall && tile != MapTileTypeBlock && tile != MapTileTypeBomb && err == nil
}

// Checks if place is empty
func (gameMap *GameMap) IsEmpty(x uint, y uint) bool {
	tile, _ := gameMap.GetTile(x, y)
	return tile == MapTileTypeEmpty
}

// Changes a tile at specified cordinates
func (gameMap *GameMap) SetTile(x uint, y uint, tile MapTile) (err error) {
	if x < MAP_MIN_X || x > MAP_MAX_X || y < MAP_MIN_Y || y > MAP_MAX_Y {
		err = ERR_INCORRECT_CORDS
		return
	}

	gameMap[(MAP_WIDTH*y)+x] = tile
	return
}

// Fetches tile at specified cordinates
func (gameMap *GameMap) GetTile(x uint, y uint) (tile MapTile, err error) {
	if x < MAP_MIN_X || x > MAP_MAX_X || y < MAP_MIN_Y || y > MAP_MAX_Y {
		err = ERR_INCORRECT_CORDS
		return
	}

	tile = gameMap[(MAP_WIDTH*y)+x]
	return
}

func createNewGameMap() (gameMap GameMap) {
	if SHOULD_GENERATE_RANDOM_MAP {
		return generateRandomMap()
	} else {
		return GameMap{ // Amazing
			MapTileTypeEmpty, MapTileTypeEmpty, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeEmpty, MapTileTypeEmpty,
			MapTileTypeEmpty, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeEmpty,
			MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock,
			MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock,
			MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock,
			MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock,
			MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock,
			MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock,
			MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock,
			MapTileTypeEmpty, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeBlock, MapTileTypeWall, MapTileTypeEmpty,
			MapTileTypeEmpty, MapTileTypeEmpty, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeBlock, MapTileTypeEmpty, MapTileTypeEmpty,
		}
	}
}

func generateRandomMap() (gameMap GameMap) {
	var col uint
	var row uint
	rand.Seed(time.Now().UnixNano())
	for col = 0; col <= MAP_MAX_X; col++ {
		for row = 0; row <= MAP_MAX_Y; row++ {
			var tile MapTile
			// corners should be empty
			if col == 0 && (row <= 1 || row >= MAP_MAX_Y-1) ||
				col == MAP_MAX_X && (row <= 1 || row >= MAP_MAX_Y-1) ||
				row == 0 && (col <= 1 || col >= MAP_MAX_X-1) ||
				row == MAP_MAX_Y && (col <= 1 || col >= MAP_MAX_X-1) {
				tile = MapTileTypeEmpty
			} else if row%2 == 1 && col%2 == 1 && col != MAP_MAX_X && row != MAP_MAX_Y {
				tile = MapTileTypeWall
			} else {
				num := rand.Intn(5)
				switch num {
				case 0, 1:
					tile = MapTileTypeEmpty
				default:
					tile = MapTileTypeBlock
				}
			}
			gameMap.SetTile(col, row, tile)
		}
	}
	return gameMap
}
