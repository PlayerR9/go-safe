package queue

// queue_node represents a node in a linked queue.
type queue_node[T any] struct {
	// value is the value stored in the node.
	value T

	// next is a pointer to the next queueLinkedNode in the queue.
	next *queue_node[T]
}
