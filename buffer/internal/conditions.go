package internal

// BufferCondition is an enumeration of the possible conditions of the Buffer.
type BufferCondition int

const (
	// IsEmpty indicates that the Buffer is empty.
	IsEmpty BufferCondition = iota

	// IsRunning indicates that the Buffer is running.
	IsRunning
)
