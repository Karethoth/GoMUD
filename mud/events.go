package mud

type Event interface {
  Update() error
  HasFinished()  bool
}



type FunctionEvent struct {
  server   *MUDServer
  trigger   Trigger
  function  func( *MUDServer ) error
  triggered bool
  finished  bool
}



func NewFunctionEvent(
    server   *MUDServer,
    trigger  Trigger,
    function func( *MUDServer ) error,
  ) *FunctionEvent {

  return &FunctionEvent{
    server,
    trigger,
    function,
    false,
    false,
  }
}



func (event *FunctionEvent) Update() error {
  if event.finished {
    return nil
  }

  triggered, finished := event.trigger.Triggered()

  if finished {
    event.finished = true
    return nil
  }

  if !triggered {
    if event.triggered {
      event.finished = true
    }
  } else {
    event.triggered = true
    event.function( event.server )
  }

  return nil
}



func (event *FunctionEvent) HasFinished() bool {
  if event.finished {
    return true
  }
  return false
}

