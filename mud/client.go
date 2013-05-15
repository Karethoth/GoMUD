package mud

import (
  "fmt"
  "net"
  "time"
  "container/list"
)


type Client struct {
  conn net.Conn

  incoming chan string
  outgoing chan string

  connected bool
  quit chan bool

  clientList *list.List

  server *MUDServer

  name string

  gameState *GameState

  activeTime time.Time
}



func NewClient( conn net.Conn, server *MUDServer, clientList *list.List ) *Client {
  newClient := &Client { 
    conn, 
    make(chan string),
    make(chan string),
    true,
    make(chan bool),
    clientList,
    server,
    "",
    nil,
    time.Now(),
  }

  // Generate name for the quest
  newClient.name = fmt.Sprintf( "Guest%s", conn.RemoteAddr() )

  // Set game state
  newClient.gameState = &GameState {
    "WelcomeScreen",
    "/welcomeScreen/start",
  }

  // Call the motd command on behalf of the player
  trigger := NewTimedTrigger( time.Now() )
  function := func( server *MUDServer ) error {
    server.games[newClient.gameState.gameName].ExecuteCommand(
      newClient,
      "motd",
    )
    return nil
  }
  server.events.PushBack( NewFunctionEvent( trigger, function ) )

  // Refresh active time
  newClient.RefreshActiveTime()

  return newClient
}



func (c *Client) Read( buffer []byte ) (bool, int) {
  count, err := c.conn.Read( buffer )
  if err != nil {
    c.Close()
    return false, 0
  }
  return true, count;
}



func (c *Client) Close() {
  c.quit <- true
  c.connected = false
  c.conn.Close()
  c.RemoveFromList()
}



func (c *Client) Equal( other *Client ) bool {
  if c.conn == other.conn {
    return true
  }
  return false
}



func (c *Client) RemoveFromList() {
  for entry := c.clientList.Front(); entry != nil; entry = entry.Next() {
    client := entry.Value.(*Client)
    if c.Equal( client ) {
      c.clientList.Remove( entry )
    }
  }
}



func (c *Client) RefreshActiveTime() {
  // Update active time
  c.activeTime = time.Now()

  // Generate timeout event
  triggerFunction := func() time.Time {
    return c.activeTime
  }

  trigger := NewTimeoutTrigger(
    c.activeTime.Add( time.Duration(500)*time.Second ),
    triggerFunction,
  )

  // Create function that's called when the timeout has triggered.
  eventFunction := func( server *MUDServer ) error {
    if !c.connected || !server.HasClient( c ) {
      return nil
    }

    fmt.Printf( "%s timed out.\n", c.name )
    c.outgoing <- "\r\n\r\nYou have timed out. Have a nice day!\r\n"
    c.Close()
    return nil
  }

  event := NewFunctionEvent( trigger, eventFunction )

  // Push the event to the event list.
  c.server.events.PushBack( event )
}



func ClientReader( client *Client, server *MUDServer ) {
  buffer := make( []byte, 2048 )
  lineBuffer := make( []byte, 2048 )
  index := 0


  for {
    ok, received := client.Read( buffer )
    if !ok {
      break
    }

    client.RefreshActiveTime()

    for i := 0; i < received; i++ {

      // If we have a line break handle the command
      // and reset lineBuffer
      if buffer[i] == '\n' {
        command := string( lineBuffer[0:index] )
        for x := 0; x < index+1; x++ {
          lineBuffer[x] = 0x00
        }
        index = 0
        
        fmt.Printf( "%s gave command '%s'.\n", client.name, command )
        if game, ok := server.games[client.gameState.gameName]; ok {
          err := game.ExecuteCommand( client, command )
          if err != nil {
            fmt.Printf( "Received error from execute command: %s\n", err.Error() )
          }
        } else {
          fmt.Printf( "game(%s) was not found\n", client.gameState.gameName )
        }

        continue
      
      // Ignore \r
      } else if buffer[i] == '\r' {
        continue

      // Ignore DEL
      } else if buffer[i] == 0x7F {
        continue

      // Handle backspace
      } else if buffer[i] == 0x8 {
        if index <= 0 {
          continue
        }

        lineBuffer[index-1] = 0x00
        index--
        continue
      }

      lineBuffer[index] = buffer[i]
      index++
    }


    for i := 0; i < 2048; i++ {
      buffer[i] = 0x00
    }
  }

  client.quit <- true
}



func ClientSender( client *Client ) {
  for {
    select {
      case buffer := <-client.outgoing:
        count := 0
        for i := 0; i < len( buffer ); i++ {
          if buffer[i] == 0x00 {
            break
          }
          count++
        }
        client.conn.Write( []byte( buffer )[0:count] )

      case <-client.quit:
        client.conn.Close()
        fmt.Printf( "%s disconnected.\n", client.name )
        return
    }
  }
}

