package graph

import (
	"testing"
)

// TestNewGraph tests the allocation of memory for a new graph
func TestNewGraph(t *testing.T) {
	g, err := NewGraph(true)

	if err != nil {
		t.Errorf("error when initializing graph %s", err)
	}

	if g.Directed != true {
		t.Errorf("expected Directed to be true, got %t", g.Directed)
	}

	if g.NVertices != 0 {
		t.Errorf("expected NVertices to be 0, got %d", g.NVertices)
	}
}

// TestReadFromFile tests the construction of a graph, using the "skeina_graph" format.
// Note that this is also indirectly testing InsertEdge, which is called by ReadFromFile
func TestReadFromFile(t *testing.T) {
	g, err := ReadFromFile("./test_example_1.skeina_graph", true)
	if err != nil {
		t.Error("error when success expected")
	}

	if g.NVertices != 8 {
		t.Errorf("expected 8 vertices, got %d", g.NVertices)
	}

	if g.NEdges != 8 {
		t.Errorf("expected 8 edges, got %d", g.NEdges)
	}

	if g.Directed != true {
		t.Error("expected directed to be true, got false")
	}
}

// TestBfs tests the breadth-first algorithm
func TestBfs(t *testing.T) {
	g, _ := ReadFromFile("./fig_5_9.skeina_graph", true)
	nVertices, nEdges := g.Bfs(1)

	// TODO: I think it's working, but unsure of better assertion criteria
	if nVertices != 6 {
		t.Errorf("Expected nVertices to be 6, got %d", nVertices)
	}
	if nEdges != 7 {
		t.Errorf("Expected nEdges to be 7, got %d", nEdges)
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
		t.Errorf("Queue should empty, but reported was reported as nonempty")
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
