package main

import (
    "bufio"
    "fmt"
    "log"
    "math/rand"
    "os"
    "time"

    "github.com/MaxHalford/gago"
    "github.com/MaxHalford/eaopt"
)

type Chromosome []int

func (c Chromosome) Evaluate(labyrinth []string) float64 {

    // Define the start and end positions
    startX, startY := 0, 0
    endX, endY := 0, 0
    for i, row := range labyrinth {
        for j, cell := range row {
            if cell == 3 {
                startX, startY = i, j
            } else if cell == 4 {
                endX, endY = i, j
            }
        }
    }

    // Traverse the labyrinth according to the chromosome and calculate the distance to the exit
    x, y := startX, startY
    for _, move := range c {
        switch move {
        case 0:
            y -= 1
        case 1:
            y += 1
        case 2:
            x -= 1
        case 3:
            x += 1
        }
        if x < 0 || x >= len(labyrinth) || y < 0 || y >= len(labyrinth[x]) || labyrinth[x][y] == 1 {
            return -1.0 // hit a green cell, invalid path
        }
        if x == endX && y == endY {
            return float64(len(c)) // reached the exit, return the path length as fitness score
        }
    }
    return -1.0 // did not reach the exit, invalid path
}

func readLabyrinth(filename string) []string {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return lines
}

func writeToFile(filename, text string, append bool) error {
    var file *os.File
    var err error
    if append {
        file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
    } else {
        file, err = os.Create(filename)
    }
    if err != nil {
        return err
    }
    defer file.Close()

    if _, err := file.WriteString(text); err != nil {
        return err
    }

    return nil
}

func (c Chromosome) Clone() gago.Genome {
	clone := make(Chromosome, len(c))
	copy(clone, c)
	return clone
}


func NewGenome(rng *rand.Rand) gago.Genome {
    chrom := make(Chromosome, 50)
    for j := range chrom {
        chrom[j] = rng.Intn(4)
    }
    return chrom
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Read the labyrinth from a file
	labyrinth := readLabyrinth("labyrinth.txt")

	// Define the genetic algorithm parameters
	ga := gago.Generational(
        gago.NewGenome(NewGenome),
		100,                  // Population size
		0,                    // Number of elites
		0.5,                  // Crossover probability
		gago.Tournament{3},   // Selection method
		gago.Convergence{1e-6, 100}, // Stop criteria
	)

	// Initialize the population with random chromosomes
/* 	n := 50 */
/* 	pop := make(gago.Population, n) */
    pop := make(gago.Population, 100)
	for i := range pop {
		pop[i] = NewGenome(rand.New(rand.NewSource(time.Now().UnixNano())))
	}

    // Run the genetic algorithm and print the best path found
    ga.Initialize(pop)
    for i := 1; i <= 100; i++ {
        ga.Enhance()
        best := pop.Best(1)[0].Genome.(Chromosome)
        fitness := ga.Best(1)[0].Fitness
/*         best := ga.Minimize(func(genome gago.Genome) float64 {
            return genome.(Chromosome).Evaluate(labyrinth)
        })        
        fitness := pop.Best(1)[0].Fitness */
        text := fmt.Sprintf("Generation: %d\nFitness score: %f\nBest path: %v\n", i, fitness, best)
        for _, move := range best {
            switch move {
            case 0:
                text += "Left "
            case 1:
                text += "Right "
            case 2:
                text += "Up "
            case 3:
                text += "Down "
            }
        }
        text += "\n\n"
        err := writeToFile("best_path.txt", text, true)
        if err != nil {
            log.Fatal(err)
        }
    }
}