package lfucache

import (
	"errors"
	"iter"

	"github.com/igoroutine-courses/gonature.lfucache/internal/linkedlist"
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
	Put(key K, value V)
	All() iter.Seq2[K, V]
	Size() int
	Capacity() int
	GetKeyFrequency(key K) (int, error)
}

type cacheNode[K comparable, V any] struct {
	key       K
	value     V
	frequency int
}

type cacheImpl[K comparable, V any] struct {
	capacity     int
	linkedList   linkedlist.List[cacheNode[K, V]]
	freeNodes    linkedlist.List[cacheNode[K, V]]
	keyToNode    map[K]*linkedlist.Node[cacheNode[K, V]]
	freqMap      map[int]*freqBasket[K, V]
	freeBaskets  map[int]*freqBasket[K, V]
	defaultValue V
}

type freqBasket[K comparable, V any] struct {
	head *linkedlist.Node[cacheNode[K, V]]
	tail *linkedlist.Node[cacheNode[K, V]]
}

func New[K comparable, V any](capacity ...int) *cacheImpl[K, V] {
	capCache := DefaultCapacity
	if len(capacity) > 0 {
		if capacity[0] <= 0 {
			panic("Capacity must be greater than zero!")
		}
		capCache = capacity[0]
	}

	return &cacheImpl[K, V]{
		capacity:    capCache,
		linkedList:  linkedlist.New[cacheNode[K, V]](),
		freeNodes:   linkedlist.New[cacheNode[K, V]](),
		keyToNode:   make(map[K]*linkedlist.Node[cacheNode[K, V]]),
		freqMap:     make(map[int]*freqBasket[K, V]),
		freeBaskets: make(map[int]*freqBasket[K, V]),
	}
}

func (c *cacheImpl[K, V]) moveToNextBasket(node *linkedlist.Node[cacheNode[K, V]]) {
	oldFreq := node.Data.frequency
	newFreq := oldFreq + 1
	currBasket := c.freqMap[oldFreq]

	node.Data.frequency = newFreq

	if nextBasket, exist := c.freqMap[newFreq]; exist {
		switch {
		case currBasket.head == currBasket.tail:
			c.freeBaskets[oldFreq] = currBasket
			c.freeBaskets[oldFreq].head = nil
			c.freeBaskets[oldFreq].tail = nil
			delete(c.freqMap, oldFreq)
		case currBasket.head == node:
			currBasket.head = currBasket.head.Next()
		case currBasket.tail == node:
			currBasket.tail = currBasket.tail.Prev()
		}

		c.linkedList.MoveAfter(node, nextBasket.tail)
		nextBasket.tail = node
	} else {
		var nextBasket *freqBasket[K, V]

		if len(c.freeBaskets) > 0 {
			for k, v := range c.freeBaskets {
				nextBasket = v
				nextBasket.head = node
				nextBasket.tail = node
				delete(c.freeBaskets, k)
				break
			}
		} else {
			nextBasket = &freqBasket[K, V]{head: node, tail: node}
		}

		c.freqMap[newFreq] = nextBasket

		switch {
		case currBasket.head == currBasket.tail:
			c.freeBaskets[oldFreq] = currBasket
			c.freeBaskets[oldFreq].head = nil
			c.freeBaskets[oldFreq].tail = nil
			delete(c.freqMap, oldFreq)
		case currBasket.head == node:
			currBasket.head = currBasket.head.Next()
			c.linkedList.MoveAfter(node, currBasket.tail)
		case currBasket.tail == node:
			currBasket.tail = currBasket.tail.Prev()
		default:
			c.linkedList.MoveAfter(node, currBasket.tail)
		}
	}
}

func (c *cacheImpl[K, V]) Get(key K) (V, error) {
	if link, ok := c.keyToNode[key]; ok {
		c.moveToNextBasket(link)
		return link.Data.value, nil
	}
	return c.defaultValue, ErrKeyNotFound
}

func (c *cacheImpl[K, V]) displacementNode() {
	minFreq := c.linkedList.Head().Data.frequency
	boardBasket := c.freqMap[minFreq]

	nodeToRemove := boardBasket.head
	delete(c.keyToNode, nodeToRemove.Data.key)
	nodeToRemove.Data.key = *new(K)
	nodeToRemove.Data.value = *new(V)
	nodeToRemove.Data.frequency = 0

	c.freeNodes.PushFront(nodeToRemove)

	if boardBasket.head == boardBasket.tail {
		c.linkedList.RemoveNode(boardBasket.head)
		c.freeBaskets[minFreq] = boardBasket
		boardBasket.head = nil
		boardBasket.tail = nil
		delete(c.freqMap, minFreq)
	} else {
		boardBasket.head = boardBasket.head.Next()
		c.linkedList.RemoveNode(boardBasket.head.Prev())
	}
}

func (c *cacheImpl[K, V]) resetNode(node *linkedlist.Node[cacheNode[K, V]], key K, value V) {
	node.SetNext(nil)
	node.SetPrev(nil)
	node.Data.value = value
	node.Data.key = key
	node.Data.frequency = 1
}

func (c *cacheImpl[K, V]) createNewNode(key K, value V) *linkedlist.Node[cacheNode[K, V]] {
	return &linkedlist.Node[cacheNode[K, V]]{
		Data: cacheNode[K, V]{
			key:       key,
			value:     value,
			frequency: 1,
		},
	}
}

func (c *cacheImpl[K, V]) insertNewNode(key K, value V) {
	var newNode *linkedlist.Node[cacheNode[K, V]]

	if c.freeNodes.Size() > 0 {
		newNode = c.freeNodes.Head()
		c.freeNodes.RemoveNode(newNode)
		c.resetNode(newNode, key, value)
	} else {
		newNode = c.createNewNode(key, value)
	}

	freq := newNode.Data.frequency

	c.keyToNode[key] = newNode
	c.linkedList.PushFront(newNode)

	if boardBasket, exist := c.freqMap[freq]; exist {
		c.linkedList.MoveAfter(newNode, boardBasket.tail)
		boardBasket.tail = newNode
	} else {
		if len(c.freeBaskets) > 0 {
			for k, v := range c.freeBaskets {
				c.freqMap[freq] = v
				delete(c.freeBaskets, k)
				break
			}
		} else {
			c.freqMap[freq] = &freqBasket[K, V]{}
		}

		boardBasket = c.freqMap[freq]

		boardBasket.head = newNode
		boardBasket.tail = newNode
	}
}

func (c *cacheImpl[K, V]) Put(key K, value V) {
	if node, exists := c.keyToNode[key]; exists {
		node.Data.value = value
		c.moveToNextBasket(node)
		return
	}

	if c.linkedList.Size() >= c.capacity {
		c.displacementNode()
	}

	c.insertNewNode(key, value)
}

func (c *cacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		current := c.linkedList.Tail()
		for current != nil {
			if !yield(current.Data.key, current.Data.value) {
				return
			}
			current = current.Prev()
		}
	}
}

func (c *cacheImpl[K, V]) Size() int {
	return c.linkedList.Size()
}

func (c *cacheImpl[K, V]) Capacity() int {
	return c.capacity
}

func (c *cacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	if link, ok := c.keyToNode[key]; ok {
		return link.Data.frequency, nil
	}
	return 0, ErrKeyNotFound
}
