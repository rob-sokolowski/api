// Package graph. This package follows along with ch. 5 of Skiena's The Algorithm Design Manual
// The example code in the book is written in C, but I'm following along using Go
// This is for educational purposes only, do not use this code!
package graph

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
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

// begin region: graph structure
// I think this could/should be dynamically set, but I'm following the book as-is
const maxVertices = 1000

type EdgeNode struct {
	Y      int
	Weight *int
	Next   *EdgeNode
}

type Node struct {
	id  int
	val interface{}
}

type Edge struct {
	src  int
	dest int
}

type Graph2 struct {
	Nodes    []interface{}
	Edges    mapset.Set[Edge]
	Directed bool
}

func (g *Graph2) addNode(n interface{}) {
	g.Nodes = append(g.Nodes, n)
}

func (g *Graph2) addEdge(src int, dest int) {
	e := Edge{
		src:  src,
		dest: dest,
	}

	if !g.Directed {
		f := Edge{
			src:  dest,
			dest: src,
		}
		g.Edges.Add(f)
	}

	g.Edges.Add(e)
}

type Graph struct {
	Edges     [maxVertices + 1]*EdgeNode
	Degree    [maxVertices + 1]int
	NVertices int
	NEdges    int
	Directed  bool
}

func NewGraph(directed bool) (g *Graph) {
	g = new(Graph)
	g.Directed = directed
	return g
}

type NodeVal interface {
	int | string
}

// JsonGraph is a simple JSON format for describing a graph
type JsonGraph[T NodeVal] struct {
	NVertices int                `json:"nVertices"`
	Directed  bool               `json:"directed"`
	Edges     [][]int            `json:"edges"`
	Nodes     []persistedNode[T] `json:"nodes"`
}

type persistedNode[T NodeVal] struct {
	Type_ string      `json:"type"`
	Val   interface{} `json:"Val"`
}

// FromJsonFile2 reads a JSON file with the schema of the JsonGraph struct,
// into a Graph2
func FromJsonFile2[T NodeVal](path string) (*Graph2, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Now let's unmarshall the data into `payload`
	var payload JsonGraph[T]
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil, err
	}

	g := Graph2{
		Edges:    mapset.NewSet[Edge](),
		Nodes:    make([]interface{}, 0, 1000),
		Directed: payload.Directed,
	}
	// insert edge
	for _, edge := range payload.Edges {
		g.addEdge(edge[0], edge[1])
	}

	// insert Nodes
	for i, pNode := range payload.Nodes {
		node, err := parseNode[T](i, pNode)
		if err != nil {
			return nil, fmt.Errorf("unable to parse persisted node, %s", pNode)
		}
		g.addNode(node)
	}

	return &g, nil
}

func parseNode[T NodeVal](i int, pnode persistedNode[T]) (interface{}, error) {
	switch pnode.Type_ {
	case "int":
		node := Node{
			id:  i,
			val: int(pnode.Val.(float64)),
		}

		return node, nil
	case "string":
		node := Node{
			id:  i,
			val: pnode.Val.(string),
		}

		return node, nil
	}

	return nil, fmt.Errorf("unknown type_ field when parsing persisted node %s", pnode.Type_)
}

type OldJsonGraph = struct {
	NVertices int     `json:"nVertices"`
	Directed  bool    `json:"directed"`
	Edges     [][]int `json:"edges"`
}

func FromJsonFile(path string) (*Graph, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Now let's unmarshall the data into `payload`
	var payload OldJsonGraph
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil, err
	}

	g := NewGraph(payload.Directed)
	g.NVertices = payload.NVertices

	for _, edge := range payload.Edges {
		err := g.InsertEdge(edge[0], edge[1], payload.Directed)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

// InsertEdge inserts an edge into the graph. If Directed the is will be from lhs to rhs, otherwise two
// Edges will be inserted, one from lhs to rhs and one from rhs to lhs
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
		// increment edge count, note that for undirected case we do not double count
		g.NEdges++
	}

	return nil
}

// end region: graph structure

// ConnectedComponents counts the distinct maximal components of a graph.
// For example, two cliques without any mutual friends between any members would
// result in two. Two cliques where just one member from each know each other would be one.
func (g *Graph) ConnectedComponents() int {
	s := newBfsSearch(g)
	c := 0
	crawlResult := new(CrawlerCounter)

	for i := 0; i < g.NVertices; i++ {
		if s.discovered[i] == false {
			c++
			g.Bfs(i, crawlResult)
		}
	}

	return c
}

func (g *Graph) Dfs(v int, crawler Crawler) (*DfsSearch, error) {
	s := newDfsSearch(g)
	return dfs_(g, v, s, crawler)
}

// Crawler provides processing instructions for searches on a graph
// TODO: naming is hard, but sticking with Crawler for now
type Crawler interface {
	ProcessVertexEarly(int)
	ProcessVertexLate(int)
	ProcessEdge(int, int)
}

type CrawlerCounter struct {
	NVerticesProcessedEarly int
	NEdgesProcessed         int
	NVerticesProcessedLate  int
}

func (r *CrawlerCounter) ProcessVertexEarly(v int) {
	r.NVerticesProcessedEarly++
}

func (r *CrawlerCounter) ProcessVertexLate(v int) {
	r.NEdgesProcessed++
}

func (r *CrawlerCounter) ProcessEdge(v, y int) {
	r.NVerticesProcessedLate++
}

// dfs_
// TODO: Does GoLang support tail-optimization for recursive funcs??
func dfs_(g *Graph, v int, s *DfsSearch, crawler Crawler) (*DfsSearch, error) {
	var p *EdgeNode
	var y int

	s.discovered[v] = true
	s.time++
	s.entryTimes[v] = s.time

	crawler.ProcessVertexEarly(v)

	p = g.Edges[v]
	for p != nil {
		y = p.Y
		if s.discovered[y] == false {
			s.parents[y] = v
			crawler.ProcessEdge(v, y)
			dfs_(g, y, s, crawler)
		} else if (!s.processed[y] && s.parents[v] != y) || g.Directed {
			crawler.ProcessEdge(v, y)
		}

		p = p.Next
	}
	crawler.ProcessVertexLate(v)

	s.time++
	s.exitTimes[v] = s.time
	s.processed[v] = true
	return s, nil
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
func (g *Graph) Bfs(start int, crawler Crawler) (int, int) {
	var parent [maxVertices + 1]int
	var tmpNode *EdgeNode
	var v int // current vertex
	var y int // successor vertex
	edgesProcessed := 0
	verticesProcessed := 0
	queue := NewQueue[int](10)

	s := newBfsSearch(g)

	queue.Enqueue(start)
	s.discovered[start] = true

	for !queue.IsEmpty() {
		v = *queue.Dequeue()
		crawler.ProcessVertexEarly(v)
		s.processed[v] = true
		tmpNode = g.Edges[v]
		for tmpNode != nil {
			y = tmpNode.Y
			if s.processed[y] == false || g.Directed {
				crawler.ProcessEdge(v, y)
			}
			if s.discovered[y] == false {
				queue.Enqueue(y)
				s.discovered[y] = true
				parent[y] = v
			}
			tmpNode = tmpNode.Next
		}

		crawler.ProcessVertexLate(v)
	}

	return verticesProcessed, edgesProcessed
}
