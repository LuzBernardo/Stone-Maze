package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Coord struct {
	row, col int
}

type Chromosome struct {
	genes []string
	fitness int
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

    // Define queue and visited set
    queue := []struct {
        path []string
    }{{[]string{}}}
    visited := make(map[Coord]bool)
    visited[start] = true

    // Define directions
    directions := []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
    directionNames := []string{"UP", "LEFT", "RIGHT", "DOWN"}

    // Perform breadth-first search
    for len(queue) > 0 {
        front := queue[0]
        queue = queue[1:]
        if front.row == dest.row && front.col == dest.col {
            path := append(front.path, "DESTINATION")
            for i := 1; i < len(path); i++ {
                fmt.Printf("%s ", path[i])
            }
            fmt.Println()
            return
        }
        for i, d := range directions {
            newCoord := Coord{front.row + d.row, front.col + d.col}
            if newCoord.row < 0 || newCoord.row >= len(maze) || newCoord.col < 0 || newCoord.col >= len(maze[newCoord.row]) {
                continue
            }
            if maze[newCoord.row][newCoord.col] == 1 || visited[newCoord] {
                continue
            }
			newPath := append(front.path, directionNames[i])
            queue = append(queue, struct {
                row, col int
                path     []string
            }{newCoord.row, newCoord.col, newPath})
            visited[newCoord] = true
			// Print the current movement direction
			fmt.Printf("%s ", directionNames[i])
        }

        // Update the maze based on the rules
        maze = updateMaze(maze)
    }

    fmt.Println("No path found")
}
