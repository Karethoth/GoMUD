package mud

import (
  "fmt"
  "net"
  "time"
  "bytes"
  "container/list"
  "code.google.com/goconf"
)



// Server-related error message.
type MUDServerError struct {
  When time.Time
  What string
}



func (e MUDServerError) Error() string {
  return fmt.Sprintf( "%v: %v", e.When, e.What )
}




// The server itself. Before Starting it, it requires working config file.
type MUDServer struct {
  conn net.Listener
  config *conf.ConfigFile

  clientList *list.List

  // map of different games
  games map[string] Game
  events           *list.List
}



// Validates provided ConfigFile for required entries and sets it as the config.
// Returns non-nil error if errors have occured.
func (server *MUDServer) SetConfig( config *conf.ConfigFile ) error {
  if( server == nil ) {
    return MUDServerError {
      time.Now(),
      "ConfigFile provided to SetConfig is nil!",
    }
  }

  requiredFields := [...][3]string {
    { "server", "port", "int" },
    { "database", "host", "string" },
    { "database", "user", "string" },
    { "database", "password", "string" },
    { "database", "database", "string" },
  }

  // Create slice to contain errors for fields.
  fieldErrors := make( []error, len( requiredFields ) )

  // Check required fields for errors
  for _, y := range requiredFields {
    var err error

    if( y[2] == "int" ) {
      _, err = config.GetInt( y[0], y[1] )
    } else if( y[2] == "string" ) {
      _, err = config.GetString( y[0], y[1] )
    }

    if( err != nil ) {
      fieldErrors = append( fieldErrors, err )
    }
  }


  // Loop trough errors and generate error message out of them.
  if( len( fieldErrors ) > 0 ) {
    var errorCount = 0
    errorBuf := bytes.NewBufferString( "\n\tFollowing errors occured when checking fields of the configuration file:\n" )

    for _, e := range fieldErrors {
      if( e == nil ) {
        continue
      }
      fmt.Fprintf( errorBuf, "\t - %s\n", e.Error() )
      errorCount++
    }

    if( errorCount > 0 ) {
      return MUDServerError {
        time.Now(),
        errorBuf.String(),
      }
    }
  }

  server.config = config
  return nil
}



// Starts the server.
func (server *MUDServer) Start() error {
  if( server.config == nil ) {
    return MUDServerError {
      time.Now(),
      "Server startup failed. No (working) config file provided.",
    }
  }


  server.clientList = list.New()

  // Initialize games
  server.games = make( map[string] Game )
  server.games["WelcomeScreen"] = InitWelcomeScreen()

  // Create list for events
  server.events = list.New()


  // Get the port that we should start listening.
  port, _ := server.config.GetInt( "server", "port" )

  var err error

  // Open the listening socket.
  server.conn, err = net.Listen( "tcp", fmt.Sprintf( ":%d", port ) )
  if( err != nil ) {
    return MUDServerError {
      time.Now(),
      err.Error(),
    }
  }


  // Start the event handler
  go EventHandler( server )

  for {
    connection, err := server.conn.Accept()
    if err != nil {
      return MUDServerError {
        time.Now(),
        err.Error(),
      }
    } else {
      go ClientHandler( connection, server )
    }
  }

  return nil
}



// Closes the server
func (server *MUDServer) Close() {
  server.conn.Close()
}



// Returns true if client is found from client list
func (server *MUDServer) HasClient( client *Client ) bool {
  for e := server.clientList.Front(); e != nil; e = e.Next() {
    if client == e.Value.(*Client) {
      return true
    }
  }
  return false
}



// Handles the initialization of new clients
func ClientHandler( conn net.Conn, server *MUDServer ) {
  addr := conn.RemoteAddr()
  fmt.Printf( "New client connected from %s\n", addr )
  newClient := NewClient( conn, server, server.clientList )
  go ClientSender( newClient )
  go ClientReader( newClient, server )
  server.clientList.PushBack( newClient )
}



// Handles the events of the server
func EventHandler( server *MUDServer ) {
  for {
    for e := server.events.Front(); e != nil; e = e.Next() {
      if e.Value.(Event).HasFinished() {
        server.events.Remove( e )
        continue
      }

      e.Value.(Event).Update( server )
    }
    time.Sleep( 100 *time.Millisecond )
  }
}

