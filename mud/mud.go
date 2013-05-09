package mud

import (
  "fmt"
  "net"
  "time"
  "bytes"
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

  return nil
}



// Closes the server
func (server *MUDServer) Close() {
  server.conn.Close()
}

