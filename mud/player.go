package mud

import "container/list"


type Player struct {
  client *Client
  name   string

  playerList *list.List
}
