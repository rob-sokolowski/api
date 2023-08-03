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

// end region: queue implementation
// begin region: stack implementation

// Stack is your textbook FILO data structure
type Stack[T any] struct {
	Elements []*T
}

// NewStack allocates memory for a stack
func NewStack[T any](cap int) *Stack[T] {
	return &Stack[T]{
		Elements: make([]*T, 0, cap),
	}
}

// Push an element atop the stack
func (s *Stack[T]) Push(e T) {
	s.Elements = append(s.Elements, &e)
}

// Pop returns the first value on the stack, removing it from the stack
func (s *Stack[T]) Pop() *T {
	el := s.Elements[len(s.Elements)-1]
	s.Elements = s.Elements[0 : len(s.Elements)-1]
	return el
}

// Peek returns the first value on the stack, keeping it on the stack
func (s *Stack[T]) Peek() *T {
	return s.Elements[len(s.Elements)-1]
}

// IsEmpty returns true if the stack has no elements, false otherwise
func (s *Stack[T]) IsEmpty() bool {
	return len(s.Elements) == 0
}

// end region: stack implementation

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

	return g, nil
}

// GraphFromFile reads a graph from a file
// Note the book example reads from user-supplied stdin. I'm following the same formatting, but instead reading from a file
//
// The format of the file is:
// The 1st line has two integers, n and m, which are the number of vertices and edges, respectively.
// The next m lines contain the edges; each two integers, x and y, the vertices of the edge
func GraphFromFile(filename string, directed bool) (g *Graph, err error) {
	g, err = NewGraph(directed)
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

	i := 0
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

// ConnectedComponents counts the distinct maximal components of a graph.
// For example, two cliques without any mutual friends between any members would
// result in two. Two cliques where just one member from each know each other would be one.
func (g *Graph) ConnectedComponents() int {
	s := newBfsSearch(g)
	c := 0

	for i := 0; i < g.NVertices; i++ {
		if s.discovered[i] == false {
			c++
			g.Bfs(i)
		}
	}

	return c
}

func (g *Graph) Dfs(v int) (*DfsSearch, error) {
	s := newDfsSearch(g)
	return dfs_(g, v, s)
}

func dfsProcessVertexEarly(v int) {
	fmt.Printf("processVertexEarly %d", v)
}

func dfsProcessVertexLate(v int) {
	fmt.Printf("processVertexLate %d", v)
}

func dfsProcessEdge(v, y int) {
	fmt.Printf("processEdge %d, %d", v, y)
}

// dfs_
// TODO: Does GoLang support tail-optimization for recursive funcs??
func dfs_(g *Graph, v int, s *DfsSearch) (*DfsSearch, error) {
	var p *EdgeNode
	var y int

	if false {
		// TODO: Where does finished come from?
		return s, nil
	}

	s.discovered[v] = true
	s.time++
	s.entryTimes[v] = s.time

	dfsProcessVertexEarly(v)

	p = g.Edges[v]
	for p != nil {
		y = p.Y
		if s.discovered[y] == false {
			s.parents[y] = v
			dfsProcessEdge(v, y)
			dfs_(g, y, s)
		} else if (!s.processed[y] && s.parents[v] != y) || g.Directed {
			dfsProcessEdge(v, y)
		}

		if false {
			// TODO: Where does finished come from?
			return s, nil
		}

		p = p.Next
	}
	dfsProcessVertexLate(v)

	s.time++
	s.exitTimes[v] = s.time
	s.processed[v] = true
	return s, nil
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

type bfsSearchHelp struct {
	graph      *Graph
	processed  []bool
	discovered []bool
}

func newBfsSearch(g *Graph) *bfsSearchHelp {
	return &bfsSearchHelp{
		graph:      g,
		processed:  make([]bool, maxVertices+1),
		discovered: make([]bool, maxVertices+1),
	}
}

type DfsSearch struct {
	discovered []bool
	time       int
	entryTimes []int
	parents    []int
	exitTimes  []int
	processed  []bool
}

func newDfsSearch(g *Graph) *DfsSearch {
	// TODO: I think maxVertices should belong to the graph object and not this module???
	fmt.Println(g)

	return &DfsSearch{
		discovered: make([]bool, maxVertices+1),
		time:       0,
		entryTimes: make([]int, maxVertices+1),
		parents:    make([]int, maxVertices+1),
		exitTimes:  make([]int, maxVertices+1),
		processed:  make([]bool, maxVertices+1),
	}
}

// Bfs performs a breadth-first search over the Graph, starting at the node start
func (g *Graph) Bfs(start int) (int, int) {
	var parent [maxVertices + 1]int
	var tmpNode *EdgeNode
	var v int // current vertex
	var y int // successor vertex
	edgesProcessed := 0
	verticesProcessed := 0
	queue := NewQueue[int](10)

	var processEarly = func(v int) {
		// noop
	}

	var processEdge = func(v, y int) {
		edgesProcessed++
	}

	var processLate = func(y int) {
		verticesProcessed++
	}

	s := newBfsSearch(g)

	queue.Enqueue(start)
	s.discovered[start] = true

	for !queue.IsEmpty() {
		v = *queue.Dequeue()
		processEarly(v)
		s.processed[v] = true
		tmpNode = g.Edges[v]
		for tmpNode != nil {
			y = tmpNode.Y
			if s.processed[y] == false || g.Directed {
				processEdge(v, y)
			}
			if s.discovered[y] == false {
				queue.Enqueue(y)
				s.discovered[y] = true
				parent[y] = v
			}
			tmpNode = tmpNode.Next
		}

		processLate(v)
	}

	return verticesProcessed, edgesProcessed
}
