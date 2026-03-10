package event

type Reason string

const (
	ReasonCreated             Reason = "Created"
	ReasonQueued              Reason = "Queued"
	ReasonSchedulingStarted   Reason = "SchedulingStarted"
	ReasonSchedulingSucceeded Reason = "SchedulingSucceeded"
	ReasonSchedulingFailed    Reason = "SchedulingFailed"
	ReasonReserved            Reason = "Reserved"
	ReasonBound               Reason = "Bound"
	ReasonStarted             Reason = "Started"
	ReasonSucceeded           Reason = "Succeeded"
	ReasonFailed              Reason = "Failed"
	ReasonPreempted           Reason = "Preempted"
	ReasonReleased            Reason = "Released"
)
