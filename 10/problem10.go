package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "math"
)

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

    asteroids := make([]Asteroid, 0)
    y := 0
    for scanner.Scan() {
        t := scanner.Text()
        for x, v := range t {
            if v == '#'{
                asteroids = append(asteroids, Asteroid{x, y})
            }
        }
        y += 1
    }

    max_vis := 0
    for j := 0; j < len(asteroids); j++ {
        //Outer loop testing each asteroid as a command point
        t_a := asteroids[j]
        slopes := make(map[float64] *struct{}, 0)

        for i := 0; i < len(asteroids); i ++ {
            if i == j {
                // Ignore same asteroid
                continue
            }
            a := asteroids[i]
            slope := math.Atan2(float64(a.Y - t_a.Y), float64(a.X-t_a.X))
            slopes[slope] = nil
        }

        if len(slopes) > max_vis {
            max_vis = len(slopes)
        }
    }
    fmt.Println("Problem 10a:", max_vis)
}

type Asteroid struct {
    X int
    Y int
}
