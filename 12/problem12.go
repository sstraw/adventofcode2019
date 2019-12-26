package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "regexp"
    "strconv"
    "math"
)

type Moon struct {
    Pos [3]int
    Vel [3]int
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

    moons := make([]*Moon,0)
    re := regexp.MustCompile(`-?\d+`)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        m := &Moon{}
        raw_coords := re.FindAll([]byte(scanner.Text()), -1)
        for i := 0; i < 3; i++ {
            v, _ := strconv.Atoi(string(raw_coords[i]))
            m.Pos[i] = v
        }
        moons = append(moons, m)
    }

    for i := 0; i < 1000; i++ {
        for _, m := range(moons) {
            for _, m2 := range(moons) {
                m.CalculateVelocity(m2)
            }
        }
        for _, m := range(moons) {
            m.ApplyVelocity()
        }
    }

    s := 0
    for _, m := range(moons) {
        fmt.Print(m)
        e := m.CalculateEnergy()
        fmt.Printf(" %4d energy\n", e)
        s += e
    }
    fmt.Println("Problem 12a: ", s)
}

func (m *Moon) CalculateVelocity(m2 *Moon) {
    for i := 0; i < 3; i++ {
        if m.Pos[i] < m2.Pos[i] {
            m.Vel[i] += 1
        } else if m.Pos[i] > m2.Pos[i] {
            m.Vel[i] -= 1
        }
    }
}

func (m *Moon) ApplyVelocity() {
    for i := 0; i < 3; i++ {
        m.Pos[i] += m.Vel[i]
    }
}

func (m *Moon) CalculateEnergy() int {
    sp := 0
    sv := 0
    for i := 0; i < 3; i++ {
        sp += int(math.Abs(float64(m.Pos[i])))
        sv += int(math.Abs(float64(m.Vel[i])))
    }
    return sp*sv
}

func (m Moon) String() string {
    fs := "pos=<x=%3d, y=%3d, z=%3d>, vel=<x=%3d, y=%3d, z=%3d>"
    return fmt.Sprintf(fs,
        m.Pos[0],
        m.Pos[1],
        m.Pos[2],
        m.Vel[0],
        m.Vel[1],
        m.Vel[2])
}

