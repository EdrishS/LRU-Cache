package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const SIZE = 5

type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Cache struct {
	Queue Queue
	Hash  Hash
}

func (c *Cache) Get(key string) (string, bool) {
	val, ok := c.Hash[key]
	if ok {
		c.Check(key)
	}
	return val.Val, ok
}

type Hash map[string]*Node

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func (c *Cache) Set(key, value string) {
	node := &Node{}

	if val, ok := c.Hash[key]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: key}
	}

	node.Val = value
	c.Add(node)
	c.Hash[key] = node
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}

	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}

func (c *Cache) Check(str string) {
	node := &Node{}

	if val, ok := c.Hash[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: str}
	}
	c.Add(node)
	c.Hash[str] = node
}

func (c *Cache) Remove(n *Node) *Node {
	fmt.Printf("remove: %s\n", n.Val)
	left := n.Left
	right := n.Right

	left.Right = right
	right.Left = left
	c.Queue.Length -= 1
	delete(c.Hash, n.Val)
	return n
}

func (c *Cache) Add(n *Node) {
	fmt.Printf("add: %s\n", n.Val)
	tmp := c.Queue.Head.Right

	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = tmp
	tmp.Left = n

	c.Queue.Length++
	if c.Queue.Length > SIZE {
		c.Remove(c.Queue.Tail.Left)
	}
}

func (c *Cache) Display() {
	c.Queue.Display()
}

func (q *Queue) Display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Length-1 {
			fmt.Printf("<-->")
		}
		node = node.Right
	}
	fmt.Printf("]")
}

func main() {
	fmt.Println("START CACHE SERVER")
	cache := NewCache()

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")

		var jsonData map[string]string
		err := json.NewDecoder(r.Body).Decode(&jsonData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cache.Set(jsonData["key"], jsonData["value"])
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		key := r.URL.Query().Get("key")
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		val, ok := cache.Get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Key not found")
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, val)
		}
	})

	http.ListenAndServe(":8080", nil)
}
