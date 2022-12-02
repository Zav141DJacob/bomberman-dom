package main

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type WebsocketMessageType uint

type WebsocketMessage struct {
	Type    WebsocketMessageType `json:"type"`
	Content any                  `json:"content"`
}

// In case you don't want to send any additional info
type WebsocketMessageSmall struct {
	Type WebsocketMessageType `json:"type"`
}

func (w WebsocketMessage) Construct() ([]byte, error) {
	if w.Content == nil {
		return json.Marshal(WebsocketMessageSmall{
			Type: w.Type,
		})
	}

	return json.Marshal(w)
}

const (
	WebsocketMessageTypeSetUsername   WebsocketMessageType = iota + 1 // Sets username
	WebsocketMessageTypeGlobalMessage                                 // Global chat message

	WebsocketMessageTypeSetOwnID   // When user joins server sends them their own ID
	WebsocketMessageTypeDisconnect // User disconnects
	WebsocketMessageTypeConnect    // User connects

	WebsocketMessageTypeJoinLobby      // User requests to join lobby
	WebsocketMessageTypeLobbyMessage   // Lobby message
	WebsocketMessageTypeGameStart      // If there are 4 members users can request to start game
	WebsocketMessageTypeGameStartTimer // Notifies users about game starting early
	WebsocketMessageTypeAnnounceWinner // sent to users when game ends (1 player alive)

	WebsocketMessageTypePlayerMove          // received and sent for user movement
	WebsocketMessageTypeThrottledPlayerMove // Sent to user when PlayerMove is called too fast
	WebsocketMessageTypeBombPlant           // sent to server when user plants a bomb
	WebsocketMessageTypeBombExplode         // sent to users when bomb explodes

	WebsocketMessageTypeItemDrop   // sent to user when something drops after breaking block
	WebsocketMessageTypeItemPickup // sent to users when somebody picks stuff up

	WebsocketMessageTypeSendGameMap // sends initial map to players
	WebsocketMessageTypeLoseLife    // sent to users when a player loses a life (from explosion for example)
)

type SessionInfo struct {
	ID, Username string
	LobbyID      uuid.UUID

	IsDead           bool      // If the player is dead
	Lives            uint      // Amount of lives the player has left
	Bombs            uint      // Amount of bombs in inventory
	MaxAmountOfBombs uint      // Amount of bombs you can hold in your inventory
	LastMoveTime     time.Time // Used for throttle

	FlamePowerup uint
	SpeedPowerup uint

	X uint // Player X position
	Y uint // Player Y position
}

type BombInfo struct {
	LobbyID uuid.UUID

	X            uint // Bomb X position
	Y            uint // Bomb Y position
	UserID       string
	FlamePowerup uint
}

type Lobby struct {
	Map            GameMap
	Members        []*melody.Session
	GameHasStarted bool
	GameHasEnded   bool

	PlantedBombs   []BombInfo
	GameStartTimer StartTimer
}

type PlayerMoveDirectionType uint

const (
	PlayerMoveDirectionTypeUp = iota + 1
	PlayerMoveDirectionTypeDown
	PlayerMoveDirectionTypeRight
	PlayerMoveDirectionTypeLeft
)

type ChatMessage struct {
	ID, Username, Content, Time string
}

// Handle user connecting to websocket
func (app *Application) websocketConnectHandler(s *melody.Session) {
	id := uuid.NewString()
	app.logDebug.Printf("%s connects to the websocket\n", id)
	s.Set(WEBSOCKET_USERINFO_LOCATION, &SessionInfo{
		ID:               id,
		Username:         "",
		LobbyID:          uuid.Nil,
		IsDead:           false,
		Lives:            DEFAULT_LIVES,
		Bombs:            DEFAULT_BOMBS,
		MaxAmountOfBombs: DEFAULT_BOMBS,
		LastMoveTime:     time.Now().Add(-(NORMAL_THROTTLE + time.Second)), // just in case

		X: 0,
		Y: 0,
	})

	raw, err := WebsocketMessage{
		Type: WebsocketMessageTypeSetOwnID,
		Content: UUIDResponse{
			ID: id,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, id)
		return
	}

	rawConnectEvent, err := WebsocketMessage{
		Type: WebsocketMessageTypeConnect,
		Content: UUIDResponse{
			ID: id,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, id)
		return
	}

	if err = s.Write(raw); err != nil { // Broadcasts own ID to session
		app.logError.Println(err)
	}

	if err = app.m.BroadcastOthers(rawConnectEvent, s); err != nil { // Broadcasts new user joining to user others
		app.logError.Println(err)
	}
}

func (app *Application) websocketMessageHandler(s *melody.Session, rawMessage []byte) {
	var message WebsocketMessage
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	// Current user info
	userInfo, exists := GetInfo(s)
	if !exists {
		return
	}

	switch message.Type {
	case WebsocketMessageTypeSetUsername:
		app.HandleSetUsername(s, userInfo, message, rawMessage)
	case WebsocketMessageTypeGlobalMessage:
		app.HandleGlobalMessage(s, userInfo, message, rawMessage)
	case WebsocketMessageTypeJoinLobby:
		app.HandleJoinLobby(s, userInfo, message, rawMessage)
	case WebsocketMessageTypeLobbyMessage:
		app.HandleLobbyMessage(s, userInfo, message, rawMessage)
	case WebsocketMessageTypeGameStart:
		app.HandleGameStart(userInfo)
	case WebsocketMessageTypePlayerMove:
		app.HandlePlayerMove(s, userInfo, message, rawMessage)
	case WebsocketMessageTypeBombPlant:
		app.HandleBombPlant(s, userInfo, message, rawMessage)
	default:
		app.logWarning.Printf("Unknown event: %d\n", message.Type)
	}
}

// Handle user disconnecting from websocket
func (app *Application) websocketDisconnectHandler(s *melody.Session) {
	userInfo, exists := GetInfo(s)
	if !exists {
		return
	}

	if userInfo.LobbyID != uuid.Nil { // If user in a lobby
		currentLobby := app.Lobbies[userInfo.LobbyID]

		// -1 to know if it cant find current user in lobby
		foundPos := -1
		for j, v := range currentLobby.Members {
			tempUserInfo, exists := GetInfo(v)
			if !exists {
				continue
			}

			if tempUserInfo.ID == userInfo.ID {
				foundPos = j
				break
			}
		}

		if foundPos != -1 { // if found current user in lobby
			currentLobby.Members = ArrayRemoveByIndex(foundPos, currentLobby.Members)

			if len(currentLobby.Members) == 0 { // If all members left, delete the lobby
				app.logDebug.Printf("Deleting lobby %s as its empty!\n", userInfo.LobbyID.String())
				currentLobby.StopTimer()
				delete(app.Lobbies, userInfo.LobbyID)
			} else {
				currentLobby.RestartTimer()
				if currentLobby.GameHasStarted && !currentLobby.GameHasEnded {
					app.HandleAnnounceWinner(s, userInfo.LobbyID)
				}
			}
		}
	}

	app.logDebug.Printf("%s disconnects from the websocket\n", userInfo.ID)

	raw, err := WebsocketMessage{
		Type: WebsocketMessageTypeDisconnect,
		Content: UUIDResponse{
			ID: userInfo.ID,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, userInfo.ID)
		return
	}

	if err = app.m.BroadcastOthers(raw, s); err != nil { // Notifies other users of this session leaving
		app.logError.Println(err)
	}
}
