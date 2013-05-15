package mud

import "time"


// Triggered() returns two bools, which have following meanings:
// - Is the trigger triggered.
// - Has the trigger went "bad" or "old".
type Trigger interface {
  Triggered() (bool, bool)
}



// Triggers when specified time is reached.
type TimedTrigger struct {
  targetTime time.Time
  triggered  bool
}



func NewTimedTrigger( targetTime time.Time ) *TimedTrigger {
  return &TimedTrigger {
    targetTime,
    false,
  }
}



func (trigger *TimedTrigger) Triggered() (bool, bool) {
  if trigger.triggered {
    return false, false
  }

  timeNow := time.Now()
  if trigger.targetTime.Before( timeNow ) ||
     trigger.targetTime == timeNow {
    trigger.triggered = true
    return true, false
  }

  return false, false
}




// Triggers when timeout is reached.
// Needs function to fetch newest "last active time".
type TimeoutTrigger struct {
  timeoutTime time.Time
  activeTime  time.Time
  GetNewTime  func() time.Time

  triggered   bool
}



func NewTimeoutTrigger(
    timeoutTime time.Time,
    GetNewTime  func() time.Time ) *TimeoutTrigger {


  return &TimeoutTrigger{
    timeoutTime,
    GetNewTime(),
    GetNewTime,
    false,
  }
}




func (trigger *TimeoutTrigger) Triggered() (bool, bool) {
  if trigger.triggered {
    return false, false
  }

  newTime := trigger.GetNewTime()

  // Check if the trigger has gotten bad.
  if newTime != trigger.activeTime {
    return false, true
  }

  // Check if the timeout time has been reached.
  if time.Now().After( trigger.timeoutTime ) {
    return true, false
  }

  // Trigger has not gone bad, but isn't active yet.
  return false, false
}

