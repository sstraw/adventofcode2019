package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strings"
    "strconv"
    "math"
)

type Coord struct {
    X int
    Y int
    D int
}

func main() {
    fname := "input.txt"
    if (len(os.Args)) == 2 {
        fname = os.Args[1]
    }
    file, err := os.Open(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    scanner.Scan()
    wire1 := strings.Split(scanner.Text(), ",")
    scanner.Scan()
    wire2 := strings.Split(scanner.Text(), ",")

    wiremap1 := MapWire(wire1)
    fmt.Printf("Wiremap1: %v\n", len(wiremap1))
    wiremap2 := MapWire(wire2)
    fmt.Printf("Wiremap2: %v\n", len(wiremap2))

    distance    := 0
    distance_3b := 0
    for _, w1 := range(wiremap1){
        for _, w2 := range(wiremap2){
            if (w1.X == w2.X &&
                w1.Y == w2.Y){
                    t_distance := ManhattanDistance(w1.X, w1.Y)
                    if (distance == 0 || t_distance < distance) {
                        distance = t_distance
                    }
                    if (distance_3b == 0 || (w1.D + w2.D) < distance_3b) {
                        distance_3b = w1.D + w2.D
                    }
            }
        }
    }
    fmt.Printf("Problem 3a: %v\n", distance)
    fmt.Printf("Problem 3b: %v\n", distance_3b)
}

func MapWire(path []string) []Coord {
    coords := make([]Coord, 0)
    x, y, d:= 0, 0, 0
    for _, instruct := range(path){
        dx, dy    := 0, 0
        steps, _  := strconv.Atoi(instruct[1:])
        switch d := instruct[0]; d {
            case 'U':
                dy = 1
            case 'D':
                dy = -1
            case 'R':
                dx = 1
            case 'L':
                dx = -1
            default:
                log.Fatal("Why the fuck are we here")
        }
        for i:=0; i < steps; i++ {
            x += dx
            y += dy
            d += 1
            coords = append(coords, Coord{x,y, d})
        }
    }
    return coords
}

func ManhattanDistance(x, y int) int {
    return int(math.Abs(float64(x)) + math.Abs(float64(y)))
}
