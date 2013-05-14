package mud

import (
  "fmt"
  "time"
  "container/list"
)


type WelcomeScreen struct {
  
}



func (game WelcomeScreen) GetAvailableCommands( client *Client ) *list.List {
  commands := list.New()
  commands.PushBack( "quit" )
  return commands
}



func (game WelcomeScreen) ExecuteCommand( client *Client, command string ) error {
  if command == "quit" {
    client.outgoing <- "Good bye!\r\n"
    client.Close()
    return nil
  }

  return MUDServerError {
    time.Now(),
    fmt.Sprintf( "Command '%s' is not a valid command.", command ),
  }
}



func InitWelcomeScreen() WelcomeScreen {
  welcomeScreen := WelcomeScreen {

  }

  return welcomeScreen
}
