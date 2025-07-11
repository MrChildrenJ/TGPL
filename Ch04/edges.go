// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 99.

// Graph shows how to use a map of maps to represent a directed graph.
package main

import "fmt"

//!+
var graph = make(map[string]map[string]bool)//graph: {str: edges}	edges: {str:bool}

func addEdge(from, to string) {				// "a" "b"							"a" "d"
	edges := graph[from]					// edges = graph["a"] -> nil		edges := graph["a"] = {"b":true}
	if edges == nil {						// O								X
		edges = make(map[string]bool)		// edges = {}		
		graph[from] = edges					// graph["a"] = edges ({})
	}
	edges[to] = true						// edges = {"b":true}				edges["d"]=true -> {"b":true, "d":true}
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

//!-

func main() {
	addEdge("a", "b")
	addEdge("c", "d")
	addEdge("a", "d")
	addEdge("d", "a")
	fmt.Println(hasEdge("a", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("a", "d"))
	fmt.Println(hasEdge("d", "a"))
	fmt.Println(hasEdge("x", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("x", "d"))
	fmt.Println(hasEdge("d", "x"))

	fmt.Println(graph)
}