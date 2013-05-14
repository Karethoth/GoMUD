package mud

import (
  "container/list"
)


type GameState struct {
  gameName string
  roomName string
}



type Game interface {
  GetAvailableCommands( *Client ) *list.List
  ExecuteCommand( *Client, string ) error
}

