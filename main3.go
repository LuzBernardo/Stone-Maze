package main

import (
    "container/heap"
    "fmt"
    "io/ioutil"
    "math"
    "os"
    "strings"
)

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
    item := x.(*Node)
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[0 : n-1]
    return item
}


type Coord struct {
	row, col int
}

type Node struct {
    coord Coord
    path  []string
    cost  float64
}

// Updates the maze based on the green and white cell rules
func updateMaze(maze [][]int) [][]int {
    numRows, numCols := len(maze), len(maze[0])
    newMaze := make([][]int, numRows)
    for i := range newMaze {
        newMaze[i] = make([]int, numCols)
        for j := range newMaze[i] {
            newMaze[i][j] = maze[i][j]
        }
    }
    for i := range maze {
        for j := range maze[i] {
            greenCount := 0
            for r := i - 1; r <= i+1; r++ {
                for c := j - 1; c <= j+1; c++ {
                    if r == i && c == j {
                        continue
                    }
                    if r < 0 || r >= numRows || c < 0 || c >= numCols {
                        continue
                    }
                    if maze[r][c] == 1 {
                        greenCount++
                    }
                }
            }
            if maze[i][j] == 0 && greenCount > 1 && greenCount < 5 {
                newMaze[i][j] = 1
            } else if maze[i][j] == 1 && (greenCount < 4 || greenCount > 5) {
                newMaze[i][j] = 0
            }
        }
    }
    return newMaze
}

func euclideanDistance(a, b Coord) float64 {
    return math.Sqrt(float64((a.row-b.row)*(a.row-b.row) + (a.col-b.col)*(a.col-b.col)))
}


func main() {
    // Read the maze from file
    content, err := ioutil.ReadFile("maze.txt")
    if err != nil {
        panic(err)
    }
    lines := strings.Split(string(content), "\n")
    var maze [][]int
    for _, line := range lines {
        if line == "" {
            continue
        }
        vals := strings.Split(line, " ")
        row := make([]int, len(vals))
        for i, val := range vals {
            fmt.Sscanf(val, "%d", &row[i])
        }
        maze = append(maze, row)
    }

    // Define starting point and destination point
    var start, dest Coord
    for i := range maze {
        for j := range maze[i] {
            if maze[i][j] == 3 {
                start = Coord{i, j}
            }
            if maze[i][j] == 4 {
                dest = Coord{i, j}
            }
        }
    }

    // Add queue priority for A*
    pq := make(PriorityQueue, 0)
    heap.Init(&pq)
    heap.Push(&pq, &Node{coord: start, cost: 0})

    // Define visited set
    visited := make(map[Coord]bool)
    visited[start] = true

    // Define directions
    directions := []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
    directionNames := []string{"U", "L", "R", "D"}

	var bestPath []string
	var bestMovements int = math.MaxInt64
	var reportInterval int = 1000

    // Implementation of A*
    pathLimit := 10000000000
    pathCount := 0
    for pq.Len() > 0 {
		pathCount++
        current := heap.Pop(&pq).(*Node)
        if current.coord == dest {
            path := append(current.path, "DESTINATION")
            movements := len(current.path)
            pathStr := fmt.Sprintf("It took %d movements: %s", movements, strings.Join(path, " "))

            // Read the existing content of the best_path.txt file
            content, err := ioutil.ReadFile("best_path.txt")
            if err != nil {
                panic(err)
            }
            existingPaths := strings.Split(string(content), "\n")

            // Check if the current path is shorter than any previous paths
            isNewBest := true
            for _, line := range existingPaths {
                if line == "" {
                    continue
                }
                var existingMovements int
                _, _ = fmt.Sscanf(line, "It took %d movements:", &existingMovements)
                if movements >= existingMovements {
                    isNewBest = false
                    break
                }
            }

            if isNewBest {
                // Append the new best path to the file
                file, err := os.OpenFile("best_path.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                if err != nil {
                    panic(err)
                }
                defer file.Close()

                _, err = file.WriteString(pathStr + "\n")
                if err != nil {
                    panic(err)
                }

				bestPath = append(current.path, "DESTINATION")
				bestMovements = movements
			
				// Add this condition to print the best path periodically
				if pathCount % reportInterval == 0 {
					fmt.Printf("Best path found so far after %d paths explored:\n", pathCount)
					fmt.Printf("It took %d movements: %s\n", bestMovements, strings.Join(bestPath, " "))
				}
			} else {
				fmt.Printf("Current path is not the new best.\n%s\n", pathStr)
			}

            // Continue the search if the path limit has not been reached
            if pathCount >= pathLimit {
                break
            }
        }
        visited[current.coord] = true
        maze = updateMaze(maze)

        for i, d := range directions {
            newCoord := Coord{current.coord.row + d.row, current.coord.col + d.col}
            if newCoord.row < 0 || newCoord.row >= len(maze) || newCoord.col < 0 || newCoord.col >= len(maze[newCoord.row]) {
                continue
            }
            if maze[newCoord.row][newCoord.col] == 1 || visited[newCoord] {
                continue
            }
            newPath := append(current.path, directionNames[i])
            newNode := &Node{
                coord: newCoord,
                path:  newPath,
                cost:  current.cost + 1 + euclideanDistance(newCoord, dest),
            }
            heap.Push(&pq, newNode)
        }
    }

	if len(bestPath) > 0 {
		fmt.Printf("Best path found so far after %d paths explored:\n", pathCount)
		fmt.Printf("It took %d movements: %s\n", bestMovements, strings.Join(bestPath, " "))
	} else {
		fmt.Println("No best path found")
	}
}

