package databases

import (
	"db_meta/dbstructs"
	"log"
	"strconv"
)

// Function to find SCCs using Kosaraju's algorithm
func FindSCCs(graphResponse *dbstructs.GraphResponse) [][]string {
	// Creating a map for tracking visited nodes and a stack for storing nodes
	visited := make(map[string]bool)
	var stack []string

	// First DFS pass to fill the stack
	for _, node := range graphResponse.Nodes {
		if !visited[node.Data.Name] {
			dfs1(node.Data, graphResponse, visited, &stack)
		}
	}

	// Creating the transposed graph
	transposedGraph := transposeGraph(graphResponse)

	// Second DFS pass for SCCs
	visited = make(map[string]bool)
	var sccs [][]string
	for len(stack) > 0 {
		nodeId := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[nodeId] {
			var scc []string
			id, err := strconv.Atoi(nodeId)
			if err != nil {
				log.Println(err)
			}
			dfs2(transposedGraph.Nodes[id].Data, transposedGraph, visited, &scc)
			sccs = append(sccs, scc)
		}
	}

	log.Println(sccs)
	return sccs
}

func dfs1(node *dbstructs.NodeData, graph *dbstructs.GraphResponse, visited map[string]bool, stack *[]string) {
	visited[node.Name] = true
	for _, edge := range graph.Edges {
		if edge.Data.Source == node.Name && !visited[edge.Data.Target] {
			dfs1(findNodeData(edge.Data.Target, graph), graph, visited, stack)
		}
	}
	*stack = append(*stack, node.Name)
}

func dfs2(node *dbstructs.NodeData, graph *dbstructs.GraphResponse, visited map[string]bool, scc *[]string) {
	visited[node.Name] = true
	*scc = append(*scc, node.Name)
	for _, edge := range graph.Edges {
		if edge.Data.Target == node.Name && !visited[edge.Data.Source] {
			dfs2(findNodeData(edge.Data.Source, graph), graph, visited, scc)
		}
	}
}

func transposeGraph(originalGraph *dbstructs.GraphResponse) *dbstructs.GraphResponse {
	transposedGraph := &dbstructs.GraphResponse{
		Nodes: make([]*dbstructs.NodeElement, len(originalGraph.Nodes)),
		Edges: []*dbstructs.RelationshipEdge{},
	}

	// Copying nodes into the transposed graph
	for i, node := range originalGraph.Nodes {
		transposedGraph.Nodes[i] = &dbstructs.NodeElement{
			Data: &dbstructs.NodeData{
				Name:       node.Data.Name,
				ID:         node.Data.ID,
				Columns:    node.Data.Columns, // check in case smthn isn't right
				PrimaryKey: node.Data.PrimaryKey,
				Indexes:    node.Data.Indexes,
			},
		}
	}

	// Inverting the edges
	for _, edge := range originalGraph.Edges {
		transposedEdge := &dbstructs.RelationshipEdge{
			Data: &dbstructs.EdgeData{
				ID:     edge.Data.ID,
				Source: edge.Data.Target,
				Target: edge.Data.Source,
			},
		}
		transposedGraph.Edges = append(transposedGraph.Edges, transposedEdge)
	}

	return transposedGraph
}

func findNodeData(nodeId string, graph *dbstructs.GraphResponse) *dbstructs.NodeData {
	for _, node := range graph.Nodes {
		if node.Data.Name == nodeId {
			return node.Data
		}
	}
	return nil
}
