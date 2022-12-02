package main

// UUIDResponse is for when you need to only respond with ID
type UUIDResponse struct {
	ID string `json:"ID"`
}

// CordsResponse is for when you need to only respond with cords
type CordsResponse struct {
	X uint `json:"X"`
	Y uint `json:"Y"`
}

// UserInfoResponse is used for notifying username change and user joining lobby
type UserInfoResponse struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
}

type PlayerMoveResponse struct {
	ID        string                  `json:"ID"`
	X         uint                    `json:"X"`
	Y         uint                    `json:"Y"`
	Direction PlayerMoveDirectionType `json:"Direction"`
}

type ItemDropResponse struct {
	Item MapTile `json:"Item"`
	X    uint    `json:"X"`
	Y    uint    `json:"Y"`
}

type ItemPickupResponse struct {
	UserID string  `json:"UserID"`
	Item   MapTile `json:"Item"`
	X      uint    `json:"X"`
	Y      uint    `json:"Y"`
}

type LoseLifeResponse struct {
	ID        string `json:"ID"`
	LivesLeft uint   `json:"LivesLeft"`
}

type Block struct {
	X uint `json:"X"`
	Y uint `json:"Y"`
}

type BombPlant struct {
	X      uint   `json:"X"`
	Y      uint   `json:"Y"`
	UserID string `json:"UserID"`
}

type BombExplodeResponse struct {
	X              uint    `json:"X"`
	Y              uint    `json:"Y"`
	UserID         string  `json:"UserID"`
	IsFlamable     bool    `json:"IsFlamable"`
	AffectedBlocks []Block `json:"AffectedBlocks"`
}

type ThrottledPlayerMoveResponse struct {
	SecondsLeft float64 `json:"SecondsLeft"`
}

type GameStartTimerResponse struct {
	Timer float64 `json:"Timer"`
}
