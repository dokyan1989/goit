package todo

type Todo struct {
	ID     int64
	Title  string
	Status string
}

const (
	TodoStatusUnknown = "UNKNOWN"
	TodoStatusNew     = "NEW"
	TodoStatusDoing   = "DOING"
	TodoStatusPending = "PENDING"
	TodoStatusDone    = "DONE"
)
