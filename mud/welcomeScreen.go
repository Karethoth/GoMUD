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

  // If client wants to quit
  if command == "quit" {
    // Prepare the event
    trigger := NewTimedTrigger( time.Now() )
    function := func ( server *MUDServer ) error {
      client.outgoing <- "Good bye!\r\n"
      client.Close()
      return nil
    }

    // Add the event to list
    client.server.events.PushBack( NewFunctionEvent( trigger, function ) )

  // MOTD
  } else if command == "motd" {
    SendMOTD( client )

  // No command found
  } else {
    trigger := NewTimedTrigger( time.Now() )
    function := func ( server *MUDServer ) error {
      client.outgoing <- fmt.Sprintf( "Command '%s' is not a valid command.\r\n> ", command )
      return nil
    }
    client.server.events.PushBack( NewFunctionEvent( trigger, function ) )
  }

  return nil
}



func InitWelcomeScreen() WelcomeScreen {
  welcomeScreen := WelcomeScreen {

  }

  return welcomeScreen
}



func SendMOTD( client *Client ) {
  message := fmt.Sprintf(
               "%c[1mWelcome!%c[0m\r\nThis server isn't useful for anything, thus I recommend you to '%c[31mquit%c[0m'.\r\n> ",
               0x1b, 0x1b, 0x1b, 0x1b,
             )

  client.outgoing <- message
}
