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

    moons := make([]Moon,0)
    re := regexp.MustCompile(`-?\d+`)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        m := Moon{}
        raw_coords := re.FindAll([]byte(scanner.Text()), -1)
        for i := 0; i < 3; i++ {
            v, _ := strconv.Atoi(string(raw_coords[i]))
            m.Pos[i] = v
        }
        moons = append(moons, m)
    }

    initial := make([]Moon, len(moons))
    copy (initial, moons)
    for i := 0; i < 1000; i++ {
        for j, _ := range(moons) {
            for k, _ := range(moons) {
                moons[j].CalculateVelocity(moons[k])
            }
        }
        for j, _ := range(moons) {
            moons[j].ApplyVelocity()
        }
        if i == 999 {
        }
    }
    s := 0
    for _, m := range(moons) {
        fmt.Print(m)
        e := m.CalculateEnergy()
        fmt.Printf(" %4d energy\n", e)
        s += e
    }
    fmt.Println("Problem 12a:", s)

    // Determine how often each repeats
    var repeats [3]int
    copy (moons, initial)
    for i := 1; ; i++ {
        for j, _ := range(moons) {
            for k, _ := range(moons) {
                moons[j].CalculateVelocity(moons[k])
            }
        }
        for j, _ := range(moons) {
            moons[j].ApplyVelocity()
        }
        // Check if any are back to the start
        for j := 0; j < 3; j ++ {
            eq := true
            for k, _ := range(moons) {
                if (moons[k].Pos[j] != initial[k].Pos[j] ||
                    moons[k].Vel[j] != initial[k].Vel[j]) {
                    eq = false
                    break
                }
            }
            if eq && repeats[j] == 0 {
                repeats[j] = i
            }
        }
        eq := true
        for _, v := range(repeats){
            if v == 0 {
                eq = false
            }
        }
        if eq {
            break
        }
    }

    // Find the LCM
    steps := LCM(LCM(repeats[0], repeats[1]), repeats[2])
    fmt.Println("Problem 12b", steps)
}

func (m *Moon) CalculateVelocity(m2 Moon) {
    for i := 0; i < 3; i++ {
        if m.Pos[i] < m2.Pos[i] {
            m.Vel[i] += 1
        } else if m.Pos[i] > m2.Pos[i] {
            m.Vel[i] -= 1
        }
    }
}

func LCM(a, b int) int {
    return (a * b)/GCD(a, b)
}

func GCD(a, b int) int {
    for ; b != 0; {
        a, b = b, a % b
    }
    return a
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

