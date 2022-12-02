package main

import "time"

// Application
const DEFAULT_PORT = 8080
const WEBSOCKET_USERINFO_LOCATION = "info"

// Bomberman
const BOMB_WAIT_TIME = time.Second * 5
const EXPLOSION_RADIUS uint = 1
const DEFAULT_BOMBS = 1
const DEFAULT_LIVES = 3
const POWERUP_DURATION = time.Second * 5

// Game & lobby
const MINIMUM_MEMBERS_TO_START_GAME = 2
const MAX_LOBBY_MEMBERS = 4
const INITIAL_GAME_WAIT = time.Second * 20
const SECOND_GAME_WAIT = time.Second * 10

// Player move throttle
const NORMAL_THROTTLE = time.Millisecond * 500
const SPEED_POWERUP_THROTTLE_ADJUSTMENT = time.Millisecond * 100
const MIN_THROTTLE = time.Millisecond * 50

// Map
const MAP_WIDTH = 13
const MAP_HEIGHT = 11
const SHOULD_GENERATE_RANDOM_MAP = true

const MAP_MAX_X = MAP_WIDTH - 1
const MAP_MAX_Y = MAP_HEIGHT - 1
const MAP_MIN_X = 0
const MAP_MIN_Y = 0

// Players
const DEFAULT_PLAYER1_X = MAP_MIN_X
const DEFAULT_PLAYER1_Y = MAP_MAX_Y
const DEFAULT_PLAYER2_X = MAP_MAX_X
const DEFAULT_PLAYER2_Y = MAP_MIN_Y
const DEFAULT_PLAYER3_X = MAP_MAX_X
const DEFAULT_PLAYER3_Y = MAP_MAX_Y
const DEFAULT_PLAYER4_X = MAP_MIN_X
const DEFAULT_PLAYER4_Y = MAP_MIN_Y
