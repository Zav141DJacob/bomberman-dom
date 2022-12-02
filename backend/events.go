package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

func (app *Application) HandleSetUsername(s *melody.Session, userInfo *SessionInfo, message WebsocketMessage, rawMessage []byte) {
	switch v := message.Content.(type) {
	case string:
		userInfo.Username = v
	default:
		app.logDebug.Println("Invalid type content was sent to set username")
		return
	}

	app.logDebug.Printf("%s changes username to '%s'\n", userInfo.ID, userInfo.Username)

	raw, err := WebsocketMessage{
		Type: WebsocketMessageTypeSetUsername,
		Content: UserInfoResponse{
			ID:       userInfo.ID,
			Username: userInfo.Username,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	s.Set(WEBSOCKET_USERINFO_LOCATION, userInfo)
	if err = app.m.BroadcastOthers(raw, s); err != nil { // Notifies other users of username change
		app.logError.Println(err)
	}
}

func (app *Application) HandleGlobalMessage(s *melody.Session, userInfo *SessionInfo, message WebsocketMessage, rawMessage []byte) {
	app.logDebug.Printf("%s sends a global message\n", userInfo.ID)

	var content string
	switch v := message.Content.(type) {
	case string:
		content = v
	default:
		app.logDebug.Println("Invalid type content was sent to global message")
		return
	}

	raw, err := WebsocketMessage{
		Type: WebsocketMessageTypeGlobalMessage,
		Content: ChatMessage{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Content:  content,
			Time:     time.Now().String(),
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	if err = app.m.Broadcast(raw); err != nil { // Sends message to everyone
		app.logError.Println(err)
	}
}

func (app *Application) HandleJoinLobby(s *melody.Session, userInfo *SessionInfo, message WebsocketMessage, rawMessage []byte) {
	if userInfo.LobbyID != uuid.Nil { // If in lobby already
		app.logDebug.Printf("%s wants to join lobby but is in one already\n", userInfo.ID)
		return
	}

	app.logDebug.Printf("%s requests to join a lobby\n", userInfo.ID)

	// Finds lobby to join or creates one
	var foundLobbyID uuid.UUID
	for id, lobby := range app.Lobbies {
		if len(lobby.Members) < MAX_LOBBY_MEMBERS && !lobby.GameHasStarted {
			foundLobbyID = id
			break
		}
	}

	// None found
	if foundLobbyID == uuid.Nil {
		app.logDebug.Println("No matching lobbies found, creating one")

		// Make sure lobby with same id doesn't exist
		foundLobbyID = uuid.New()
		for _, exists := app.Lobbies[foundLobbyID]; exists; foundLobbyID = uuid.New() {
		}

		app.Lobbies[foundLobbyID] = &Lobby{
			Map:            createNewGameMap(),
			Members:        []*melody.Session{s},
			GameHasStarted: false,
			GameHasEnded:   false,
		}
	} else {
		currentLobby := app.Lobbies[foundLobbyID]
		currentLobby.Members = append(currentLobby.Members, s)
		go currentLobby.StartTimer(app, foundLobbyID)
	}

	// Sets new-found lobby ID to database
	userInfo.LobbyID = foundLobbyID
	s.Set(WEBSOCKET_USERINFO_LOCATION, userInfo)
	app.logDebug.Printf("%s is in lobby %s now! :D\n", userInfo.ID, userInfo.LobbyID)

	rawJoin, err := WebsocketMessage{
		Type: WebsocketMessageTypeJoinLobby,
		Content: UserInfoResponse{
			ID:       userInfo.ID,
			Username: userInfo.Username,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	// Sends lobby join event to all lobby members
	// and at the same time info about other already in lobby members to you
	if err = app.m.BroadcastFilter(rawJoin, func(session *melody.Session) bool {
		tempUserInfo, exists := GetInfo(session)
		if !exists || tempUserInfo.LobbyID != userInfo.LobbyID {
			return false
		}

		// Sends info about other members to you
		if tempUserInfo.ID != userInfo.ID {
			if raw, err := (WebsocketMessage{
				Type: WebsocketMessageTypeJoinLobby,
				Content: UserInfoResponse{
					ID:       tempUserInfo.ID,
					Username: tempUserInfo.Username,
				},
			}.Construct()); err != nil {
				app.logWarning.Println(err, string(rawMessage))
			} else {
				if err = s.Write(raw); err != nil { // Sends info about other users to you as well
					app.logError.Println(err)
				}
			}
		}

		return true
	}); err != nil {
		app.logError.Println(err)
	}
}

func (app *Application) HandleLobbyMessage(s *melody.Session, userInfo *SessionInfo, message WebsocketMessage, rawMessage []byte) {
	if userInfo.LobbyID == uuid.Nil { // If not in lobby
		app.logDebug.Printf("%s tried to send lobby message but isn't in a lobby\n", userInfo.ID)
		return
	}

	app.logDebug.Printf("%s sends a lobby message\n", userInfo.ID)

	var content string
	switch v := message.Content.(type) {
	case string:
		content = v
	default:
		app.logDebug.Println("Invalid type content was sent to lobby message")
		return
	}

	raw, err := WebsocketMessage{
		Type: WebsocketMessageTypeLobbyMessage,
		Content: ChatMessage{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Content:  content,
			Time:     time.Now().String(),
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	if err = app.m.BroadcastFilter(raw, BroadcastAllLobbyMembers(userInfo.LobbyID)); err != nil {
		app.logError.Println(err)
	}
}

func (app *Application) HandleGameStart(userInfo *SessionInfo) {
	if userInfo.LobbyID == uuid.Nil { // If not in lobby
		app.logDebug.Printf("%s requests to start a lobby ...but they arent in a lobby themselves :/\n", userInfo.ID)
		return
	} else if app.Lobbies[userInfo.LobbyID].GameHasStarted { // rejects game start when game already started
		app.logDebug.Printf("%s tried to start game but its already started\n", userInfo.ID)
		return
	}

	// Gets lobby member count to make sure it has 2 or more members
	currentLobbyMemberCount := len(app.Lobbies[userInfo.LobbyID].Members)
	if currentLobbyMemberCount < MINIMUM_MEMBERS_TO_START_GAME {
		app.logDebug.Printf("%s requests to start the lobby but has %d members :(\n", userInfo.ID, currentLobbyMemberCount)
		return
	}

	app.logDebug.Printf("%s starts lobby with %d members!\n", userInfo.ID, currentLobbyMemberCount)
	app.HandleGameStartInternal(userInfo.LobbyID)
}

func (app *Application) HandleGameStartInternal(LobbyID uuid.UUID) {
	rawGameStart, err := WebsocketMessage{
		Type: WebsocketMessageTypeGameStart,
	}.Construct()
	if err != nil {
		app.logWarning.Println(err)
		return
	}
	app.Lobbies[LobbyID].GameHasStarted = true

	// Sends game start event to all lobby users
	if err = app.m.BroadcastFilter(rawGameStart, BroadcastAllLobbyMembers(LobbyID)); err != nil {
		app.logError.Println(err)
	}

	// Sends initial map to players
	rawSendMap, err := WebsocketMessage{
		Type:    WebsocketMessageTypeSendGameMap,
		Content: app.Lobbies[LobbyID].Map,
	}.Construct()
	if err != nil {
		app.logWarning.Println(err)
		return
	}

	if err = app.m.BroadcastFilter(rawSendMap, BroadcastAllLobbyMembers(LobbyID)); err != nil {
		app.logError.Println(err)
	}

	// Assigns different starting positions to players
	sessions, err := app.m.Sessions()
	if err != nil {
		app.logError.Println(err)
		return
	}

	playerNum := 1
	for _, session := range sessions {
		tempUserInfo, exists := GetInfo(session)
		if !exists || tempUserInfo.LobbyID != LobbyID {
			continue
		}

		switch playerNum {
		case 1: // top left
			tempUserInfo.X = DEFAULT_PLAYER1_X
			tempUserInfo.Y = DEFAULT_PLAYER1_Y
		case 2: // bottom right
			tempUserInfo.X = DEFAULT_PLAYER2_X
			tempUserInfo.Y = DEFAULT_PLAYER2_Y
		case 3: // top right
			tempUserInfo.X = DEFAULT_PLAYER3_X
			tempUserInfo.Y = DEFAULT_PLAYER3_Y
		case 4: // bottom left
			tempUserInfo.X = DEFAULT_PLAYER4_X
			tempUserInfo.Y = DEFAULT_PLAYER4_Y
		}

		session.Set("info", tempUserInfo)
		playerNum++

		// broadcast cords to other here
		rawPlayerPos, err := WebsocketMessage{
			Type: WebsocketMessageTypePlayerMove,
			Content: PlayerMoveResponse{
				ID: tempUserInfo.ID,
				X:  tempUserInfo.X,
				Y:  tempUserInfo.Y,
			},
		}.Construct()
		if err != nil {
			app.logWarning.Println(err)
			continue
		}

		if err = app.m.BroadcastFilter(rawPlayerPos, BroadcastAllLobbyMembers(LobbyID)); err != nil {
			app.logError.Println(err)
		}
	}
}

func (app *Application) HandlePlayerMove(s *melody.Session, userInfo *SessionInfo, message WebsocketMessage, rawMessage []byte) {
	if userInfo.LobbyID == uuid.Nil { // If not in lobby
		app.logDebug.Printf("%s tried to move but they arent in a lobby\n", userInfo.ID)
		return
	} else if !app.Lobbies[userInfo.LobbyID].GameHasStarted { // If game hasn't started yet
		app.logDebug.Printf("%s tried to move but their game hasn't started yet!\n", userInfo.ID)
		return
	} else if app.Lobbies[userInfo.LobbyID].GameHasEnded { // If game is over
		app.logDebug.Printf("%s tried to move but the game has ended already\n", userInfo.ID)
		return
	} else if userInfo.IsDead { // If user is dead
		app.logDebug.Printf("%s tried to move but they are dead!\n", userInfo.ID)
		return
	}

	var moveDirection PlayerMoveDirectionType
	switch direction := message.Content.(type) {
	case float64:
		moveDirection = PlayerMoveDirectionType(direction)
	default:
		app.logDebug.Println("Invalid type content was sent to player move")
		return
	}

	timeSince := time.Since(userInfo.LastMoveTime)
	personalThrottle := NORMAL_THROTTLE - time.Duration(userInfo.SpeedPowerup)*SPEED_POWERUP_THROTTLE_ADJUSTMENT
	personalThrottle = maxDuration(personalThrottle, MIN_THROTTLE)

	if timeSince < personalThrottle {
		app.logDebug.Printf("%s tried to move but they are currently throttled", userInfo.ID)

		rawThrottledMove, err := WebsocketMessage{
			Type: WebsocketMessageTypeThrottledPlayerMove,
			Content: ThrottledPlayerMoveResponse{
				SecondsLeft: (personalThrottle - timeSince).Seconds(),
			},
		}.Construct()
		if err != nil {
			app.logWarning.Println(err, string(rawMessage))
			return
		}

		s.Write(rawThrottledMove)
		return
	}

	gamemap := app.Lobbies[userInfo.LobbyID].Map
	switch moveDirection {
	case PlayerMoveDirectionTypeUp:
		if userInfo.Y <= MAP_MIN_Y {
			app.logDebug.Printf("%s player tried to move up but cant go any higher\n", userInfo.ID)
			return
		} else if !gamemap.CanMoveTo(userInfo.X, userInfo.Y-1) {
			app.logDebug.Printf("%s player tried to move up but there is a block\n", userInfo.ID)
			return
		}

		app.logDebug.Printf("%s player moves up\n", userInfo.ID)
		userInfo.Y--
	case PlayerMoveDirectionTypeDown:
		if userInfo.Y >= MAP_MAX_Y {
			app.logDebug.Printf("%s player tried to move down but cant go any lower\n", userInfo.ID)
			return
		} else if !gamemap.CanMoveTo(userInfo.X, userInfo.Y+1) {
			app.logDebug.Printf("%s player tried to move down but there is a block\n", userInfo.ID)
			return
		}

		app.logDebug.Printf("%s player moves down\n", userInfo.ID)
		userInfo.Y++
	case PlayerMoveDirectionTypeRight:
		if userInfo.X >= MAP_MAX_X {
			app.logDebug.Printf("%s player tried to move right but cant go any more right\n", userInfo.ID)
			return
		} else if !gamemap.CanMoveTo(userInfo.X+1, userInfo.Y) {
			app.logDebug.Printf("%s player tried to move right but there is a block\n", userInfo.ID)
			return
		}

		app.logDebug.Printf("%s player moves right\n", userInfo.ID)
		userInfo.X++
	case PlayerMoveDirectionTypeLeft:
		if userInfo.X <= MAP_MIN_X {
			app.logDebug.Printf("%s player tried to move left but cant go any more left\n", userInfo.ID)
			return
		} else if !gamemap.CanMoveTo(userInfo.X-1, userInfo.Y) {
			app.logDebug.Printf("%s player tried to move left but there is a block\n", userInfo.ID)
			return
		}

		app.logDebug.Printf("%s player moves left\n", userInfo.ID)
		userInfo.X--
	default:
		app.logDebug.Printf("%s player gave invalid move direction\n", userInfo.ID)
		return
	}

	userInfo.LastMoveTime = time.Now()
	s.Set(WEBSOCKET_USERINFO_LOCATION, userInfo)

	// broadcast cords to other users
	rawPlayerPos, err := WebsocketMessage{
		Type: WebsocketMessageTypePlayerMove,
		Content: PlayerMoveResponse{
			ID:        userInfo.ID,
			X:         userInfo.X,
			Y:         userInfo.Y,
			Direction: moveDirection,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	if err = app.m.BroadcastFilter(rawPlayerPos, BroadcastAllLobbyMembers(userInfo.LobbyID)); err != nil {
		app.logError.Println(err)
	}

	app.HandlePowerupPickup(userInfo)
}

func (app *Application) HandleBombPlant(s *melody.Session, userInfo *SessionInfo, message WebsocketMessage, rawMessage []byte) {
	if userInfo.LobbyID == uuid.Nil { // If not in lobby
		app.logDebug.Printf("%s tried to plant bomb but they arent in a lobby\n", userInfo.ID)
		return
	} else if !app.Lobbies[userInfo.LobbyID].GameHasStarted { // If game hasn't started yet
		app.logDebug.Printf("%s tried to plant bomb but their game hasn't started yet!\n", userInfo.ID)
		return
	} else if app.Lobbies[userInfo.LobbyID].GameHasEnded { // If game has ended
		app.logDebug.Printf("%s tried to plant bomb but the game has ended already\n", userInfo.ID)
		return
	} else if userInfo.IsDead {
		app.logDebug.Printf("%s tried to plant bomb but they are dead!\n", userInfo.ID)
		return
	} else if userInfo.Bombs < 1 { // If no bombs
		app.logDebug.Printf("%s tried to plant bomb but didn't have any!\n", userInfo.ID)
		return
	}

	bomb := BombInfo{
		LobbyID: userInfo.LobbyID,

		X:            userInfo.X,
		Y:            userInfo.Y,
		UserID:       userInfo.ID,
		FlamePowerup: userInfo.FlamePowerup,
	}

	if ArrayContains(bomb, app.Lobbies[userInfo.LobbyID].PlantedBombs) {
		app.logDebug.Printf("%s tried to plant bomb but there is already bomb in these coordinates!\n", userInfo.ID)
		return
	}

	app.logDebug.Printf("%s planted a bomb\n", userInfo.ID)

	userInfo.Bombs--
	s.Set(WEBSOCKET_USERINFO_LOCATION, userInfo)

	// Sends bomb plant event
	rawBombPlant, err := WebsocketMessage{
		Type: WebsocketMessageTypeBombPlant,
		Content: BombPlant{
			X:      userInfo.X,
			Y:      userInfo.Y,
			UserID: userInfo.ID,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	if err = app.m.BroadcastFilter(rawBombPlant, BroadcastAllLobbyMembers(userInfo.LobbyID)); err != nil {
		app.logError.Println(err)
	}

	app.Lobbies[userInfo.LobbyID].PlantedBombs = append(app.Lobbies[userInfo.LobbyID].PlantedBombs, bomb)
	app.Lobbies[userInfo.LobbyID].Map.SetTile(bomb.X, bomb.Y, MapTileTypeBomb)
	go app.HandleBombExplode(s, bomb, userInfo, rawMessage)
}

// HandleBombExplode is only called by server itself
func (app *Application) HandleBombExplode(s *melody.Session, bomb BombInfo, userInfo *SessionInfo, rawMessage []byte) {
	time.Sleep(BOMB_WAIT_TIME) // sleep
	if !ArrayContains(bomb, app.Lobbies[bomb.LobbyID].PlantedBombs) {
		rawBombExplode, err := WebsocketMessage{
			Type: WebsocketMessageTypeBombExplode,
			Content: BombExplodeResponse{
				X:      bomb.X,
				Y:      bomb.Y,
				UserID: bomb.UserID,
			},
		}.Construct()
		if err != nil {
			app.logWarning.Println(err, string(rawMessage))
			return
		}

		if err = app.m.BroadcastFilter(rawBombExplode, BroadcastAllLobbyMembers(userInfo.LobbyID)); err != nil {
			app.logError.Println(err)
		}
		app.logDebug.Printf("Bomb at X%v Y%v already exploded!\n", bomb.X, bomb.Y)
		userInfo.Bombs++
		return
	}
	blocks, potentialPowerups := app.Lobbies[bomb.LobbyID].CalculateDestroyedBlocks(bomb)

	rawBombExplode, err := WebsocketMessage{
		Type: WebsocketMessageTypeBombExplode,
		Content: BombExplodeResponse{
			X:              bomb.X,
			Y:              bomb.Y,
			UserID:         bomb.UserID,
			IsFlamable:     bomb.FlamePowerup > 0,
			AffectedBlocks: blocks,
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err, string(rawMessage))
		return
	}

	if err = app.m.BroadcastFilter(rawBombExplode, BroadcastAllLobbyMembers(userInfo.LobbyID)); err != nil {
		app.logError.Println(err)
	}
	app.logDebug.Printf("%s's bomb exploded!\n", userInfo.ID)
	userInfo.Bombs++

	for _, blk := range potentialPowerups {
		app.HandlePowerup(blk.X, blk.Y, bomb.LobbyID)
	}

	sessions, err := app.m.Sessions()
	if err != nil {
		app.logError.Println(err)
		return
	}

	for _, session := range sessions {
		tempUserInfo, exists := GetInfo(session)
		if !exists || tempUserInfo.LobbyID != userInfo.LobbyID {
			continue
		}

		var cachedUser SessionInfo = *tempUserInfo // users are stored instead of pointed in case this part takes too long and player has already moved away
		PlayerHitByExplosion := CoordinateInSlice(Block{X: cachedUser.X, Y: cachedUser.Y}, blocks)

		if PlayerHitByExplosion {
			tempUserInfo.Lives -= 1
			session.Set(WEBSOCKET_USERINFO_LOCATION, tempUserInfo)
			app.logDebug.Printf("%s lost his life, lives left: %v\n", tempUserInfo.ID, tempUserInfo.Lives)
			rawLoseLife, err := WebsocketMessage{
				Type: WebsocketMessageTypeLoseLife,
				Content: LoseLifeResponse{
					ID:        tempUserInfo.ID,
					LivesLeft: tempUserInfo.Lives,
				},
			}.Construct()
			if err != nil {
				app.logWarning.Println(err, string(rawMessage))
				return
			}

			// notifies all users that a player has lost a life
			app.m.BroadcastFilter(rawLoseLife, BroadcastAllLobbyMembers(userInfo.LobbyID))
		}

		if tempUserInfo.Lives < 1 {
			tempUserInfo.IsDead = true
			session.Set(WEBSOCKET_USERINFO_LOCATION, tempUserInfo)
		}
	}

	app.HandleAnnounceWinner(s, userInfo.LobbyID)
}

func (lobby *Lobby) CalculateDestroyedBlocks(bomb BombInfo) (blocks []Block, potentialPowerups []Block) {
	// Look for affected cells
	var x, y uint
	blocks = []Block{{X: bomb.X, Y: bomb.Y}}

	// Returns true when it should return early from loop
	logic := func(X_coord uint, Y_coord uint) (ReturnEarly bool) {
		tile, err := lobby.Map.GetTile(X_coord, Y_coord)
		if err != nil || tile == MapTileTypeWall {
			return true
		}

		block := Block{X: X_coord, Y: Y_coord}
		blocks = append(blocks, block)
		if tile == MapTileTypeBlock {
			potentialPowerups = append(potentialPowerups, block)
			return true
		}
		return
	}

	explosionRadius := EXPLOSION_RADIUS + bomb.FlamePowerup

	for x = 1; x <= explosionRadius; x++ {
		if logic(bomb.X+x, bomb.Y) {
			break
		}
	}

	for x = 1; x <= explosionRadius; x++ {
		if logic(bomb.X-x, bomb.Y) {
			break
		}
	}

	for y = 1; y <= explosionRadius; y++ {
		if logic(bomb.X, bomb.Y+y) {
			break
		}
	}

	for y = 1; y <= explosionRadius; y++ {
		if logic(bomb.X, bomb.Y-y) {
			break
		}
	}

	lobby.PlantedBombs = ArrayRemoveByValue(bomb, lobby.PlantedBombs)
	lobby.Map.SetTile(bomb.X, bomb.Y, MapTileTypeEmpty)

	for _, bmb := range lobby.PlantedBombs {
		if CoordinateInSlice(Block{X: bmb.X, Y: bmb.Y}, blocks) {
			addBlocks, addPotentialPowerups := lobby.CalculateDestroyedBlocks(bmb)
			blocks = append(blocks, addBlocks...)
			potentialPowerups = append(potentialPowerups, addPotentialPowerups...)
		}
	}
	return
}

// HandleAnnounceWinner is only called by server itself
func (app *Application) HandleAnnounceWinner(s *melody.Session, lobbyID uuid.UUID) {
	var winnerID string
	for _, member := range app.Lobbies[lobbyID].Members {
		tempUserInfo, exists := GetInfo(member)
		if !exists {
			continue
		}

		if !tempUserInfo.IsDead {
			if winnerID != "" { // If it has found another alive user already, returns
				return
			}

			winnerID = tempUserInfo.ID
		}
	}

	app.Lobbies[lobbyID].GameHasEnded = true

	app.logDebug.Printf("%s is winner of lobby %s!", winnerID, lobbyID.String())

	rawAnnounceWinner, err := WebsocketMessage{
		Type: WebsocketMessageTypeAnnounceWinner,
		Content: UUIDResponse{
			ID: winnerID,
		},
	}.Construct()
	if err != nil {
		app.logError.Println(err)
		return
	}

	if err = app.m.BroadcastFilter(rawAnnounceWinner, BroadcastAllLobbyMembers(lobbyID)); err != nil {
		app.logError.Println(err)
	}
}

// HandlePowerup is only called by server itself
func (app *Application) HandlePowerup(x uint, y uint, lobbyID uuid.UUID) {
	tempLobby := app.Lobbies[lobbyID]
	randomNum := MapTile(RandomNumber(0, 7))
	if randomNum == MapTileTypeBombPowerup || randomNum == MapTileTypeFlamePowerup || randomNum == MapTileTypeSpeedPowerup {
		app.logDebug.Printf("Block dropped a %v powerup at X%v Y%v in lobby %s\n", randomNum, x, y, lobbyID.String())
		tempLobby.Map.SetTile(x, y, randomNum)

		rawDropItem, err := WebsocketMessage{
			Type: WebsocketMessageTypeItemDrop,
			Content: ItemDropResponse{
				Item: randomNum - 2, // Converts MapTile to frontend's Item enum
				X:    x,
				Y:    y,
			},
		}.Construct()
		if err != nil {
			app.logError.Println(err)
			return
		}

		if err = app.m.BroadcastFilter(rawDropItem, BroadcastAllLobbyMembers(lobbyID)); err != nil {
			app.logError.Println(err)
			return
		}
	} else {
		tempLobby.Map.SetTile(x, y, MapTileTypeEmpty)
	}
	app.Lobbies[lobbyID] = tempLobby
}

// HandlePowerup is only called by server itself
func (app *Application) HandlePowerupPickup(user *SessionInfo) (err error) {
	tempLobby := app.Lobbies[user.LobbyID]
	gamemap := app.Lobbies[user.LobbyID].Map

	tile, err := gamemap.GetTile(user.X, user.Y)
	if err != nil {
		app.logError.Println(err)
		return
	}

	if tile > MapTileTypeBlock {
		switch tile {
		case MapTileTypeBombPowerup:
			user.MaxAmountOfBombs++
			user.Bombs++
			app.logDebug.Printf("%s picked up a Bomb powerup\n", user.ID)
		case MapTileTypeFlamePowerup:
			user.FlamePowerup++
			app.logDebug.Printf("%s picked up a Flame powerup\n", user.ID)
		case MapTileTypeSpeedPowerup:
			user.SpeedPowerup++
			app.logDebug.Printf("%s picked up a Speed powerup\n", user.ID)
		default:
			app.logDebug.Printf("Unknown powerup: %v\n", tile)
			return
		}

		gamemap.SetTile(user.X, user.Y, MapTileTypeEmpty)
		tempLobby.Map = gamemap

		rawItemPickup, err := WebsocketMessage{
			Type: WebsocketMessageTypeItemPickup,
			Content: ItemPickupResponse{
				UserID: user.ID,
				Item:   tile - 2, // Converts MapTile to frontend's Item enum
				X:      user.X,
				Y:      user.Y,
			},
		}.Construct()
		if err != nil {
			app.logError.Println(err)
			return err
		}

		if err = app.m.BroadcastFilter(rawItemPickup, BroadcastAllLobbyMembers(user.LobbyID)); err != nil {
			app.logError.Println(err)
			return err
		}
	}

	app.Lobbies[user.LobbyID] = tempLobby
	return
}
