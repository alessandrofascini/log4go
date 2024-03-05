package appenders

import (
	"fmt"
	"testing"
	"time"
)

type ComparatorTest struct {
	s             string
	expectedValue int
}

func TestComparator_Match(t *testing.T) {
	comparator := Comparator{
		pattern: "test;.log",
		sep:     ".",
	}
	tests := []ComparatorTest{
		{
			s:             "test.log",
			expectedValue: 0,
		},
		{
			s:             "test.1.log",
			expectedValue: 1,
		},
		{
			s:             "test.leg",
			expectedValue: -1,
		},
		{
			s:             "test-1.log",
			expectedValue: -1,
		},
		{
			s:             "test..log",
			expectedValue: 0,
		},
		{
			s:             "testlog",
			expectedValue: -1,
		},
	}
	for _, test := range tests {
		res := comparator.Match(test.s)
		if res != test.expectedValue {
			fmt.Printf("test failed: %v %d\n", test, res)
			t.FailNow()
		}
	}
}

func TestComparator_Replace(t *testing.T) {
	comparator := Comparator{
		pattern: "test;.log",
		sep:     ".",
	}
	tests := []struct {
		i             int
		expectedValue string
	}{
		{
			i:             1,
			expectedValue: "test.1.log",
		},
		{
			i:             3,
			expectedValue: "test.3.log",
		},
	}
	for _, test := range tests {
		res := comparator.Replace(test.i)
		if res != test.expectedValue {
			fmt.Printf("test failed: %v %s\n", test, res)
			t.FailNow()
		}
	}
}

func TestTime(t *testing.T) {
	const days = 24 * 60 * 60 * 1000
	now := time.Now()
	tenDaysAgo, _ := time.Parse("2006-01-02", "2023-05-09")
	fmt.Println(now.UnixMilli())
	fmt.Println(tenDaysAgo.UnixMilli())
	fmt.Println((now.UnixMilli() - tenDaysAgo.UnixMilli()) / days)
}

func TestLeetCode(t *testing.T) {
	eqs := [][]string{
		{"a", "b"},
		{"b", "c"},
	}
	vals := []float64{2.0, 3.0}
	queries := [][]string{
		{"a", "c"},
		{"b", "a"},
		{"a", "e"},
		{"a", "a"},
		{"x", "x"},
	}
	fmt.Println(calcEquation(eqs, vals, queries))
}

func calcEquation(equations [][]string, values []float64, queries [][]string) []float64 {
	// create graph
	graph := make(map[string]map[string]float64)
	for i, eq := range equations {
		u := eq[0]
		v := eq[1]
		if graph[u] == nil {
			graph[u] = map[string]float64{}
		}
		graph[u][v] = values[i]
		if graph[v] == nil {
			graph[v] = map[string]float64{}
		}
		graph[v][u] = 1 / values[i]
	}

	res := make([]float64, len(queries))
	for i, query := range queries {
		if graph[query[0]] == nil || graph[query[1]] == nil {
			res[i] = -1
		} else if query[0] == query[1] {
			res[i] = 1.0
		} else {
			res[i] = dfs(graph, query[0], query[1], 1, map[string]bool{})
		}
	}
	return res
}

func dfs(graph map[string]map[string]float64, start, end string, acc float64, visited map[string]bool) float64 {
	if visited[start] {
		return -1.0
	}
	visited[start] = true
	edges := graph[start]
	if edges[end] != 0 {
		return acc * edges[end]
	}
	for k, e := range edges {
		res := dfs(graph, k, end, acc*e, visited)
		if res != -1.0 {
			return res
		}
	}
	return -1.0
}
