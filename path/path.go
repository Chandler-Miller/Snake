package path

import (
	"fmt"

	"math"
)

// Node represents a node in the grid
type Node struct {
	X, Y   int
	G, H   float64
	Parent *Node
}

// Updated Node struct to support proper comparison
func (n *Node) Equal(other *Node) bool {
	return n.X == other.X && n.Y == other.Y
}

// A* search algorithm implementation
func AStarSearch(start, dest *Node, grid [][]int) []*Node {
	// initialize an open list that only contains the start node.
	// the open list holds nodes that still need to be checked
	openList := []*Node{start}
	// initialize the closed list for nodes that have already been checked
	closedList := []*Node{}

	for len(openList) > 0 {
		current := openList[0]
		currentIndex := 0
		// Find the node with the lowest total cost (G + H) in the open list
		for i, node := range openList {
			if node.G+node.H < current.G+current.H {
				current = node
				currentIndex = i
			}
		}

		if current.Equal(dest) {
			// If the current node is the destination node, reconstruct and return the path
			return reconstructPath(current)
		}

		// Swap current node with the last node in the open list
		openList[currentIndex], openList[len(openList)-1] = openList[len(openList)-1], openList[currentIndex]
		openList = openList[:len(openList)-1]
		closedList = append(closedList, current)

		// Generate neighboring nodes of the current node
		neighbors := generateNeighbors(current, grid)

		// if len(neighbors) != 0 {
		// 	panic(*neighbors[0])
		// }

		for _, neighbor := range neighbors {
			// Skip obstacles and already closed nodes
			if grid[neighbor.Y][neighbor.X] == 1 || contains(closedList, neighbor) {
				continue
			}

			// Calculate the new cost from the start node to the neighbor node
			newG := current.G + distance(current, neighbor)

			if !contains(openList, neighbor) || newG < neighbor.G {
				// If the neighbor is not in the open list or the new cost is lower than the previous cost,
				// update its parent, G score, and H score, and add it to the open list
				neighbor.Parent = current
				neighbor.G = newG
				neighbor.H = distance(neighbor, dest)
				if !contains(openList, neighbor) {
					openList = append(openList, neighbor)
				}
			}
		}
	}

	// No path found
	return nil
}

// Reconstruct the path from the destination node to the start node

func reconstructPath(node *Node) []*Node {
	path := []*Node{node}

	for node.Parent != nil {
		node = node.Parent
		path = append([]*Node{node}, path...)
	}

	return path
}

// Generate neighboring nodes of a given node
func generateNeighbors(node *Node, grid [][]int) []*Node {
	neighbors := []*Node{}

	directions := [][]int{
		{-1, 0}, // left
		{1, 0},  // right
		{0, -1}, // up
		{0, 1},  // down
	}

	for _, dir := range directions {
		x := node.X + dir[0]
		y := node.Y + dir[1]
		// Check if the neighbor is within the grid boundaries
		if x >= 0 && x < len(grid[0]) && y >= 0 && y < len(grid) {
			neighbors = append(neighbors, &Node{X: x, Y: y})
		}
	}

	return neighbors
}

// Calculate the Euclidean distance between two nodes
func distance(a, b *Node) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

// Check if a node is present in a list

func contains(list []*Node, node *Node) bool {
	for _, n := range list {
		if n.Equal(node) {
			return true
		}
	}
	return false
}

func PrintPathOnGrid(grid [][]int, start, dest *Node, path []*Node) {
	// Create a copy of the grid to avoid modifying the original grid
	updatedGrid := make([][]int, len(grid))

	for i := range grid {
		updatedGrid[i] = make([]int, len(grid[i]))
		copy(updatedGrid[i], grid[i])
	}

	// Set the start and destination nodes on the updated grid
	updatedGrid[start.Y][start.X] = 2
	updatedGrid[dest.Y][dest.X] = 3

	// Set the path nodes on the updated grid
	for _, node := range path {
		updatedGrid[node.Y][node.X] = 4
	}

	// Print the updated grid
	for _, row := range updatedGrid {
		for _, cell := range row {
			switch cell {
			case 0:
				fmt.Print("0 ")
			case 1:
				fmt.Print("X ")
			case 2:
				fmt.Print("S ")
			case 3:
				fmt.Print("D ")
			case 4:
				fmt.Print("* ")
			}
		}
		fmt.Println()
	}
}
