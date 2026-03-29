package linkedlist

type List[T any] interface {
	Size() int
	PushFront(*Node[T])
	RemoveNode(node *Node[T])
	Head() *Node[T]
	Tail() *Node[T]
	MoveAfter(currNode, prevNode *Node[T])
}

type Node[T any] struct {
	Data       T
	next, prev *Node[T]
}

func (n *Node[T]) Next() *Node[T] {
	return n.next
}

func (n *Node[T]) Prev() *Node[T] {
	return n.prev
}

func (n *Node[T]) SetNext(next *Node[T]) {
	n.next = next
}

func (n *Node[T]) SetPrev(prev *Node[T]) {
	n.prev = prev
}

type listImpl[T any] struct {
	head *Node[T]
	tail *Node[T]
	len  int
}

func New[T any]() List[T] {
	return &listImpl[T]{}
}

func (l *listImpl[T]) Head() *Node[T] {
	return l.head
}

func (l *listImpl[T]) Tail() *Node[T] {
	return l.tail
}

func (l *listImpl[T]) Size() int {
	return l.len
}

func (l *listImpl[T]) PushFront(node *Node[T]) {
	if l.head == nil {
		l.tail = node
		l.head = node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}
	l.len++
}

func (l *listImpl[T]) RemoveNode(node *Node[T]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	node.prev = nil
	node.next = nil

	l.len--
}

func (l *listImpl[T]) MoveAfter(currNode, prevNode *Node[T]) {
	if currNode.prev != nil {
		currNode.prev.next = currNode.next
	} else {
		l.head = currNode.next
	}

	if currNode.next != nil {
		currNode.next.prev = currNode.prev
	} else {
		l.tail = currNode.prev
	}

	currNode.prev = prevNode
	currNode.next = prevNode.next

	prevNode.next = currNode

	if currNode.next != nil {
		currNode.next.prev = currNode
	} else {
		l.tail = currNode
	}
}
