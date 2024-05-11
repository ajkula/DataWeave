package databases

import (
	"db_meta/dbstructs"
	"log"
)

// Function to find SCCs using Kosaraju's algorithm
func FindSCCs(graphResponse *dbstructs.GraphResponse) [][]string {
	visited := make(map[string]bool)
	var stack []string

	// First DFS pass
	for _, node := range graphResponse.Nodes {
		if !visited[node.Data.Name] {
			dfs1(node.Data, graphResponse, visited, &stack)
		}
	}

	transposedGraph := transposeGraph(graphResponse)

	// Second DFS
	visited = make(map[string]bool)
	var sccs [][]string
	for len(stack) > 0 {
		nodeName := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[nodeName] {
			var scc []string
			dfs2(findNodeByName(nodeName, transposedGraph).Data, transposedGraph, visited, &scc)
			sccs = append(sccs, scc)
		}
	}
	log.Println(sccs)

	return sccs
}

func dfs1(node *dbstructs.NodeData, graph *dbstructs.GraphResponse, visited map[string]bool, stack *[]string) {
	visited[node.Name] = true
	for _, edge := range graph.Edges {
		targetNode := findNodeData(edge.Data.Target, graph)
		if edge.Data.Source == node.Name && !visited[targetNode.Name] {
			dfs1(targetNode, graph, visited, stack)
		}
	}
	*stack = append(*stack, node.Name)
}

func dfs2(node *dbstructs.NodeData, graph *dbstructs.GraphResponse, visited map[string]bool, scc *[]string) {
	visited[node.Name] = true
	*scc = append(*scc, node.Name)
	for _, edge := range graph.Edges {
		sourceNode := findNodeData(edge.Data.Source, graph)
		if edge.Data.Target == node.Name && !visited[sourceNode.Name] {
			dfs2(sourceNode, graph, visited, scc)
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
				Columns:    node.Data.Columns,
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
	log.Println("Original edges:", originalGraph.Edges)
	log.Println("Transposed edges:", transposedGraph.Edges)

	return transposedGraph
}

func findNodeByName(nodeName string, graph *dbstructs.GraphResponse) *dbstructs.NodeElement {
	for _, node := range graph.Nodes {
		if node.Data.Name == nodeName {
			return node
		}
	}

	return nil
}

func findNodeData(nodeId string, graph *dbstructs.GraphResponse) *dbstructs.NodeData {
	for _, node := range graph.Nodes {
		if node.Data.Name == nodeId {
			return node.Data
		}
	}

	return nil
}
