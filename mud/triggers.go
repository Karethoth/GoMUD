package mud

import "time"


type Trigger interface {
  Triggered() bool
}



type TimedTrigger struct {
  targetTime time.Time
  triggered  bool
}



func (trigger *TimedTrigger) Triggered() bool {
  if trigger.triggered {
    return false
  }

  timeNow := time.Now()
  if trigger.targetTime.Before( timeNow ) ||
     trigger.targetTime == timeNow {
    trigger.triggered = true
    return true
  }

  return false
}
