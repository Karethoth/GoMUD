package mud

import (
  "time"
  "fmt"
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
  config *conf.ConfigFile
}



// Validates provided ConfigFile for required entries and sets it as the config.
// Returns non-nil error if errors have occured.
func (server *MUDServer) SetConfig( config *conf.ConfigFile ) error {
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

  return MUDServerError {
    time.Now(),
    "Server startup failed. It's not yet written.",
  }
}



func (server *MUDServer) Close() {
}

