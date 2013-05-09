package main

import (
  "fmt"
  "github.com/Karethoth/GoMUD/mud"
  "code.google.com/goconf"
)


// Starts the server up. Requires pointer to ConfigFile.
// When no errors have occured returns server and nil.
// Otherwise return nil and error.
func StartServer( config *conf.ConfigFile ) (*mud.MUDServer, error) {
  server := new( mud.MUDServer )

  err := server.SetConfig( config )
  if( err != nil ) {
    return nil, err
  }

  err = server.Start()
  if( err != nil ) {
    return nil, err
  }

  return server, nil
}



func main() {
  fmt.Printf( "Reading configuration file..\n" )

  config, err := conf.ReadConfigFile( "server.conf" )
  if( err != nil ) {
    fmt.Printf( "Received error: %s\n", err.Error() )
    return
  }

  fmt.Printf( "Configuration file has been read succesfully.\n" )
  fmt.Printf( "Starting the server..\n" )

  server, err := StartServer( config )
  if( err != nil ) {
    fmt.Printf( "Received error: %s\n", err.Error() )
    return
  }

  fmt.Printf( "Server started succesfully!\n" )

  fmt.Printf( "Closing the server now.\n" )
  server.Close()
}

