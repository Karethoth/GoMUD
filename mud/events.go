package mud

type Event interface {
  Update( *MUDServer ) error
  HasFinished()  bool
}



type FunctionEvent struct {
  trigger   Trigger
  function  func( *MUDServer ) error
  triggered bool
  finished  bool
}



func NewFunctionEvent(
    trigger  Trigger,
    function func( *MUDServer ) error,
  ) *FunctionEvent {

  return &FunctionEvent{
    trigger,
    function,
    false,
    false,
  }
}



func (event *FunctionEvent) Update( server *MUDServer ) error {
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
    event.function( server )
  }

  return nil
}



func (event *FunctionEvent) HasFinished() bool {
  if event.finished {
    return true
  }
  return false
}

