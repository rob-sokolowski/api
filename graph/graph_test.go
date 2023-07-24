package graph

import "testing"

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
		t.Error("expected 8 vertices")
	}

	if g.NEdges != 8 {
		t.Error("expected 8 edges")
	}

	if g.Directed != true {
		t.Error("expected directed to be true")
	}
}
