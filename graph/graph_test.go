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
	g, _ := ReadFromFile("./test_example_1.skeina_graph", true)

	g.bfs(1)

}

// TestQueue tests a queue data structure. Enqueue a few elements, dequeueing them should return elements in
// the same ordered. Dequeueing an empty queue returns nil
func TestQueue(t *testing.T) {
	q, _ := NewQueue[string](2)

	q.Enqueue("a")
	q.Enqueue("b")

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
}
