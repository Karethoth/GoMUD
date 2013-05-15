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
    trigger := TimedTrigger{ time.Now(), false }
    function := func ( server *MUDServer ) error {
      client.outgoing <- "Good bye!\r\n"
      client.Close()
      return nil
    }

    // Add the event to list
    client.server.events.PushBack( NewFunctionEvent( client.server, &trigger, function ) )
    return nil
  }

  trigger := TimedTrigger{ time.Now(), false }
  function := func ( server *MUDServer ) error {
    client.outgoing <- fmt.Sprintf( "Command '%s' is not a valid command.\r\n", command )
    return nil
  }
  client.server.events.PushBack( NewFunctionEvent( client.server, &trigger, function ) )

  return nil
}



func InitWelcomeScreen() WelcomeScreen {
  welcomeScreen := WelcomeScreen {

  }

  return welcomeScreen
}
