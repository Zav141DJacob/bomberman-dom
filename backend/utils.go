package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

// Session info extraction part abstracted
func GetInfo(s *melody.Session) (userInfo *SessionInfo, exists bool) {
	value, exists := s.Get("info")
	if !exists {
		return
	}

	userInfo = value.(*SessionInfo)

	return
}

func BroadcastAllLobbyMembers(lobbyID uuid.UUID) func(session *melody.Session) bool {
	return func(session *melody.Session) bool {
		tempUserInfo, exists := GetInfo(session)
		if !exists {
			return false
		}

		return tempUserInfo.LobbyID == lobbyID
	}
}

func Distance(a uint, b uint) uint {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}

// Check if slice of blocks includes specific coordinate
func CoordinateInSlice(coord Block, blocks []Block) bool {
	for _, block := range blocks {
		if block == coord {
			return true
		}
	}
	return false
}

func ArrayRemoveByIndex[T any](index int, array []T) []T {
	array[index] = array[len(array)-1]
	return array[:len(array)-1]
}

func ArrayRemoveByValue[T comparable](value T, array []T) (newArray []T) {
	for _, v := range array {
		if v != value {
			newArray = append(newArray, v)
		}
	}
	return
}

func ArrayContains[T comparable](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// RandomNumber generates a random number in range
func RandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func maxDuration(a, b time.Duration) time.Duration {
	if a >= b {
		return a
	}
	return b
}
