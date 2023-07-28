// Package graph. This package follows along with ch. 5 of Skiena's The Algorithm Design Manual
// The example code in the book is written in C, but I'm following along using Go
// This is for educational purposes only, do not use this code!
package graph

import (
	"bufio"
	"fmt"
	"os"
)

// begin region: queue implementation

// Queue is your textbook FIFO data structure, borrowed from Ch. 3, but I made it a bit more generic
type Queue[T any] struct {
	Elements []*T
}

// NewQueue allocates memory for a queue
func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		Elements: make([]*T, 0, cap),
	}
}

func (q *Queue[T]) Dequeue() *T {
	if len(q.Elements) == 0 {
		return nil
	}

	el := q.Elements[0]
	q.Elements = q.Elements[1:]
	return el
}

func (q *Queue[T]) Enqueue(e T) {
	q.Elements = append(q.Elements, &e)
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.Elements) == 0
}

// end regions: queue implementation

// I think this could/should be dynamically set, but I'm following the book as-is
const maxVertices = 1000

type EdgeNode struct {
	Y      int
	Weight *int
	Next   *EdgeNode
}

type Graph struct {
	Edges     [maxVertices + 1]*EdgeNode
	Degree    [maxVertices + 1]int
	NVertices int
	NEdges    int
	Directed  bool
}

func NewGraph(directed bool) (g *Graph, err error) {
	g = new(Graph)

	g.NVertices = 0
	g.NEdges = 0
	g.Directed = directed

	for i := 0; i < maxVertices; i++ {
		g.Degree[i] = 0
		g.Edges[i] = nil
	}

	return g, nil
}

// ReadFromFile reads a graph from a file
// Note the book example reads from user-supplied stdin. I'm following the same formatting, but instead reading from a file
//
// The format of the file is:
// The 1st line has two integers, n and m, which are the number of vertices and edges, respectively.
// The next m lines contain the edges; each two integers, x and y, the vertices of the edge
func ReadFromFile(filename string, directed bool) (g *Graph, err error) {
	g, err = NewGraph(true)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var i int = 0
	for scanner.Scan() {
		var lhs, rhs int
		line := scanner.Text()
		// First line sets vertex / edge count, the rest are edges
		if i == 0 {
			_, err = fmt.Sscanf(line, "%d %d", &g.NVertices, &g.NEdges)
		} else {
			_, err = fmt.Sscanf(line, "%d %d", &lhs, &rhs)
			err = g.InsertEdge(lhs, rhs, directed)
		}

		if err != nil {
			return nil, err
		}
		i++
	}

	return g, nil
}

// InsertEdge inserts an edge into the graph. If directed the is will be from lhs to rhs, otherwise two
// edges will be inserted, one from lhs to rhs and one from rhs to lhs
func (g *Graph) InsertEdge(lhs, rhs int, directed bool) (err error) {
	p := new(EdgeNode)

	p.Weight = nil
	p.Y = rhs
	p.Next = g.Edges[lhs]

	g.Edges[lhs] = p
	g.Degree[lhs]++

	if directed == false {
		err = g.InsertEdge(rhs, lhs, true)
	} else {
		// TODO: I'm not sure about this.. why are directed edges doubly counted?
		//g.NEdges++
	}

	return nil
}

// Bfs performs a breadth-first search over the Graph, starting at the node start
func (g *Graph) Bfs(start int) (int, int) {
	var processed [maxVertices + 1]bool
	var discovered [maxVertices + 1]bool
	var parent [maxVertices + 1]int
	var tmpNode *EdgeNode
	var v int // current vertex
	var y int // successor vertex
	edgesProcessed := 0
	verticesProcessed := 0
	queue := NewQueue[int](10)

	var initializeSearch = func() {
		for i := 0; i < maxVertices+1; i++ {
			parent[i] = -1
		}
	}

	var processEarly = func(v int) {
		// noop
	}

	var processEdge = func(v, y int) {
		edgesProcessed++
	}

	var processLate = func(y int) {
		verticesProcessed++
	}

	initializeSearch()
	queue.Enqueue(start)
	discovered[start] = true

	for !queue.IsEmpty() {
		v = *queue.Dequeue()
		processEarly(v)
		processed[v] = true
		tmpNode = g.Edges[v]
		for tmpNode != nil {
			y = tmpNode.Y
			if processed[y] == false || g.Directed {
				processEdge(v, y)
			}
			if discovered[y] == false {
				queue.Enqueue(y)
				discovered[y] = true
				parent[y] = v
			}
			tmpNode = tmpNode.Next
		}

		processLate(v)
	}

	return verticesProcessed, edgesProcessed
}
