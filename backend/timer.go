package main

import (
	"time"

	"github.com/google/uuid"
)

type StartTimer struct {
	StartTime    time.Time
	SentReminder bool
	IsRunning    bool
	Abort        bool
}

func (lobby *Lobby) RestartTimer() {
	lobby.GameStartTimer.StartTime = time.Now()
	lobby.GameStartTimer.SentReminder = false
}

func (lobby *Lobby) StopTimer() {
	lobby.GameStartTimer.Abort = true
}

// StartTimer SHOULD ONLY BE USED IN A GOROUTINE
func (lobby *Lobby) StartTimer(app *Application, lobbyID uuid.UUID) {
	if lobby.GameHasStarted {
		return
	} else if len(lobby.Members) < MINIMUM_MEMBERS_TO_START_GAME {
		return
	} else if lobby.GameStartTimer.IsRunning {
		lobby.RestartTimer()
		return
	}

	lobby.GameStartTimer.IsRunning = true
	lobby.RestartTimer()
	defer func() {
		lobby.GameStartTimer.IsRunning = false
	}()

	rawReminder, err := WebsocketMessage{
		Type: WebsocketMessageTypeGameStartTimer,
		Content: GameStartTimerResponse{
			Timer: SECOND_GAME_WAIT.Seconds(),
		},
	}.Construct()
	if err != nil {
		app.logWarning.Println(err)
		return
	}

	for {
		if lobby.GameStartTimer.Abort || lobby.GameHasStarted { // Should abort in these cases
			return
		}

		timePassed := time.Since(lobby.GameStartTimer.StartTime)
		maxLobbyMembersMode := len(lobby.Members) >= MAX_LOBBY_MEMBERS // Different behaviour when there is max amount of lobby members

		if len(lobby.Members) >= MINIMUM_MEMBERS_TO_START_GAME {
			if !lobby.GameStartTimer.SentReminder && (maxLobbyMembersMode || timePassed >= INITIAL_GAME_WAIT) {
				app.m.BroadcastFilter(rawReminder, BroadcastAllLobbyMembers(lobbyID))
				lobby.GameStartTimer.SentReminder = true

				if maxLobbyMembersMode {
					app.logDebug.Printf(
						"Sending a game start reminder to lobby %s becuase it has max amount of members\n",
						lobbyID.String(),
					)
				} else {
					app.logDebug.Printf(
						"Sending a game start reminder to lobby %s becuase %f seconds has passed\n",
						lobbyID.String(),
						timePassed.Seconds(),
					)
				}
			} else if (maxLobbyMembersMode && timePassed >= SECOND_GAME_WAIT) || timePassed >= INITIAL_GAME_WAIT + SECOND_GAME_WAIT {
				if maxLobbyMembersMode {
					app.logDebug.Printf(
						"Starting lobby %s automatically because %f seconds has passed\n",
						lobbyID.String(),
						timePassed.Seconds(),
					)
				} else {
					app.logDebug.Printf(
						"Starting lobby %s automatically because %f seconds has passed\n",
						lobbyID.String(),
						timePassed.Seconds(),
					)
				}

				app.HandleGameStartInternal(lobbyID)
				return
			}
		}

		time.Sleep(time.Second) // Tries every seconds
	}
}
