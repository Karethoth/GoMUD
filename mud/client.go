package mud

import (
  "fmt"
  "net"
  "container/list"
)


type Client struct {
  conn net.Conn
  player *Player

  incoming chan string
  outgoing chan string

  quit chan bool

  clientList *list.List
}



func NewClient( conn net.Conn, clientList *list.List ) *Client {
  return &Client{ conn, nil, make(chan string), make(chan string), make(chan bool), clientList }
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
    client := entry.Value.(Client)
    if c.Equal( &client ) {
      c.clientList.Remove( entry )
    }
  }
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


    for i := 0; i < received; i++ {

      // If we have a line break handle the command
      // and reset lineBuffer
      if buffer[i] == '\n' {
        command := string( lineBuffer[0:index+1] )
        for x := 0; x < index+1; x++ {
          lineBuffer[x] = 0x00
        }
        index = 0
        
        fmt.Printf( "Command was: '%s'.\n", command )
        continue
      
      // Ignore \r
      } else if buffer[i] == '\r' {
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
        fmt.Println( "Client Disconnected" )
        break
    }
  }
}

