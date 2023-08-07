package graph

import (
	"fmt"
	"testing"
)

// TestNewGraph tests the allocation of memory for a new graph
func TestNewGraph(t *testing.T) {
	g := NewGraph(true)

	if g.Directed != true {
		t.Errorf("expected Directed to be true, got %t", g.Directed)
	}

	if g.NVertices != 0 {
		t.Errorf("expected NVertices to be 0, got %d", g.NVertices)
	}
}

type MyUnion interface {
	Person | Movie
}

func readVal(i int) interface{} {
	if i == 1 {
		return Person{FirstName: "Rob", Lastname: "Sokolowski"}
	} else if i == 2 {
		return Movie{Title: "Star Wars", Year: 1977}
	}
	return nil
}

type MyStruct[T MyUnion] struct {
	Vals []T
}

func TestUnionType[T MyUnion](t *testing.T) {
	var arr []interface{}

	arr = make([]interface{}, 0, 10)

	arr = append(arr, readVal(1))
	arr = append(arr, readVal(2))

	s := MyStruct[T]{
		Vals: arr,
	}

	fmt.Println(s)

	fmt.Println("Hello!")
}

func TestReadFromJsonGraph2(t *testing.T) {
	g, _ := FromJsonFile2[int]("./test_example_1_generic.json")

	if g.Edges.Cardinality() != 16 {
		t.Errorf("Set of cardinality 16 expected, got %d", g.Edges.Cardinality())
	}
	if g.Nodes.Cardinality() != 8 {
		t.Errorf("Set of cardinality 8 expected, got %d", g.Nodes.Cardinality())
	}
}

func TestReadJsonGraph(t *testing.T) {
	g, _ := FromJsonFile("./test_example_1.json")

	if g.NVertices != 8 {
		t.Errorf("expected 8 vertices, got %d", g.NVertices)
	}
	if g.NEdges != 8 {
		t.Errorf("expected 8 edges, got %d", g.NEdges)
	}
	if g.Directed != false {
		t.Errorf("expected undirected, got directed")
	}
}

// TestBfs tests the breadth-first algorithm
func TestBfs(t *testing.T) {
	g, _ := FromJsonFile("./fig_5_9.json")
	crawlResult := new(CrawlerCounter)

	g.Bfs(1, crawlResult)

	// TODO: I think it's working, but unsure of better assertion criteria
	if crawlResult.NVerticesProcessedEarly != 6 {
		t.Errorf("Expected nVertices to be 6, got %d", crawlResult.NVerticesProcessedEarly)
	}
	if crawlResult.NEdgesProcessed != 6 {
		t.Errorf("Expected nEdges to be 6, got %d", crawlResult.NEdgesProcessed)
	}
}

// TestQueue tests a queue data structure. Enqueue a few elements, dequeueing them should return elements in
// the same ordered. Dequeueing an empty queue returns nil
func TestQueue(t *testing.T) {
	q := NewQueue[string](2)

	q.Enqueue("a")
	q.Enqueue("b")

	if q.IsEmpty() {
		t.Errorf("Queue should not be empty, but reported was reported as empty")
	}

	e1 := q.Dequeue()
	if *e1 != "a" {
		t.Errorf("Expected 'a', got %s", *e1)
	}

	q.Enqueue("c")
	e2 := q.Dequeue()
	if *e2 != "b" {
		t.Errorf("Expected 'b', got %s", *e2)
	}

	e3 := q.Dequeue()
	if *e3 != "c" {
		t.Errorf("Expected 'c', got %s", *e3)
	}

	e4 := q.Dequeue()
	if e4 != nil {
		t.Errorf("Expected nil pointer, got %p", e4)
	}

	if !q.IsEmpty() {
		t.Errorf("Queue should empty, but was reported as nonempty")
	}
}

// TestStack tests stack behavior. Peek() does not modify the stack, but returns the top value
// as Pop() does modify. Push() pushes an element atop the stack
func TestStack(t *testing.T) {
	s := NewStack[string](10)

	s.Push("a")
	s.Push("b")

	p1 := s.Peek()
	if *p1 != "b" {
		t.Errorf("Expected b got %s", *p1)
	}

	p2 := s.Pop()
	if *p2 != "b" {
		t.Errorf("Expected b got %s", *p2)
	}

	if s.IsEmpty() {
		t.Error("Expected non-empty, but reported empty")
	}

	p3 := s.Pop()
	if *p3 != "a" {
		t.Errorf("Expected a got %s", *p3)
	}

	if !s.IsEmpty() {
		t.Error("Expected empty, got non-empty")
	}
}

func TestDfs(t *testing.T) {
	g, _ := FromJsonFile("./fig_5_9.json")
	crawlResult := new(CrawlerCounter)
	s, _ := g.Dfs(1, crawlResult)

	// TODO: Better test criteria??
	for i := 1; i <= 6; i++ {
		if s.processed[i] == false {
			t.Errorf("Expected node %d to be processed, but wasn't", i)
		}
	}
	fmt.Println(g, s)
}
