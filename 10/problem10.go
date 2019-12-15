package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "math"
    "sort"
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

    var m_slopes map[float64][]Asteroid
    max_vis := 0
    for j := 0; j < len(asteroids); j++ {
        //Outer loop testing each asteroid as a command point
        t_a := asteroids[j]
        slopes := make(map[float64] []Asteroid, 0)

        for i := 0; i < len(asteroids); i ++ {
            if i == j {
                // Ignore same asteroid
                continue
            }
            a := asteroids[i]
            slope := math.Atan2(float64(a.X-t_a.X), float64(a.Y - t_a.Y))
            if slope > (math.Pi/2) {
                //slope = -1.5 * math.Pi + (math.Pi-slope)
            }
            slopes[slope] = append(slopes[slope], a)
        }

        if len(slopes) > max_vis {
            max_vis = len(slopes)
            m_slopes = slopes
        }
    }
    fmt.Println("Problem 10a:", max_vis)

    // List of keys
    ks := make([]float64,0)
    for k, _ := range(m_slopes) {
        ks = append(ks, k)
    }
    sort.Float64s(ks)
    blasted_roids := 0
    for i := len(ks); i > 0; i-=1 {
        target_i := ks[i-1]
        target := m_slopes[target_i]
        if len(target) > 0 {
            blasted_roids += 1
            if blasted_roids == 200 {
                fmt.Println("Problem 10b:", target[0].X*100 + target[0].Y)
                break
            } else if len(target) > 1 {
                m_slopes[target_i] = target[1:]
            } else {
                m_slopes[target_i] = target[1:]
            }
       }
    }
}

type Asteroid struct {
    X int
    Y int
}
