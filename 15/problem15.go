package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "github.com/sstraw/adventofcode2019/intcode"
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
    scanner.Scan()

    var m Map
    m.Drone.X, m.Drone.Y = 25, 25
    m.Grid[25][25].Distance = -1
    m.OxX, m.OxY = -1, -1
    m.Drone.Rc = intcode.NewComputer(scanner.Text())


    go m.Drone.Rc.Run()

    m.Explore()
    fmt.Println(m)
    fmt.Println("Problem 15a:", m.Grid[m.OxY][m.OxX].Distance)
    fmt.Println("Problem 15b:", m.Oxygenize())
}

func (m *Map) Explore () {
    distance := 0
    for distance >= 0 {
        exploring := false
        // Drone should look in all four directions
        for i := 1; i < 5; i++ {
            dx, dy := Move(i)
            t_sq   := &m.Grid[m.Drone.Y + dy][m.Drone.X + dx]
            // Skip if square is not new or a wall
            if (t_sq.Value == 1 ||
                t_sq.Distance != 0) {
                continue
            }
            // Go explore square
            m.Drone.Rc.Input <- i
            res := <-m.Drone.Rc.Output

            // Hit a wall
            if res == 0 {
                t_sq.Value = 1
                continue
            }
            // No wall, we're moving!
            exploring = true
            m.Drone.X += dx
            m.Drone.Y += dy
            // Found oxygen
            if res == 2 {
                m.OxX, m.OxY = m.Drone.X, m.Drone.Y
            }
            distance += 1
            m.Grid[m.Drone.Y][m.Drone.X].Distance = distance
            break
        }

        // If we didn't explore/move, head "home"
        if !exploring {
            d_home := 0
            shortest_distance := 99999
            for i := 1; i < 5; i++ {
                dx, dy := Move(i)
                t_sq   := &m.Grid[m.Drone.Y + dy][m.Drone.X + dx]
                if (t_sq.Distance < shortest_distance &&
                    t_sq.Distance != 0) {
                    shortest_distance = t_sq.Distance
                    d_home = i
                }
            }
            m.Drone.Rc.Input <- d_home
            out := <-m.Drone.Rc.Output
            if out != 1 && out != 2 {
                fmt.Println(m)
                log.Fatal("Return direction hit a wall:", out)
            }
            dx, dy := Move(d_home)
            m.Drone.X += dx
            m.Drone.Y += dy
            distance = shortest_distance
        }
    }
}

func (m *Map) Oxygenize () int {
    tick := 0
    m.Grid[m.OxY][m.OxX].Value = 2
    fill := true
    for ; fill; tick++ {
        // Tracks if we're still spreading
        fill = false
        squares := make([]*Square, 0)
        for y, row := range(m.Grid) {
            for x, sq := range(row) {
                if sq.Value != 2 {
                    continue
                }
                for i := 1; i < 5; i++ {
                    dx, dy := Move(i)
                    t_sq := &m.Grid[y+dy][x+dx]
                    squares = append(squares, t_sq)
                }
            }
        }
        for _, sq := range(squares) {
            if sq.Value == 0 {
                fill = true
                sq.Value = 2
            }
        }
    }
    // Subtract one because technically, the last tick we were already fully
    // oxygenized
    return tick-1
}

type Square struct {
    Value    int
    Distance int
}

type Map struct {
    Grid   [53][50]Square
    OxX, OxY       int
    Drone Drone
}

type Drone struct {
    X, Y int
    Rc   *intcode.IntcodeComputer
}

func (m Map) String() string {
    s := ""
    for y, row := range(m.Grid) {
        for x, sq := range(row) {
            if x == m.Drone.X && y == m.Drone.Y {
                s += "D"
            } else if x == m.OxX && y == m.OxY {
                s += "O"
            } else {
                switch sq.Value {
                case 0:
                    s += " "
                case 1:
                    s += "#"
                case 2:
                    s += "."
                }
            }
        }
        s += "\n"
    }
    return s
}

// A simple function to return what the x/y
// increment should be givent the direction
func Move(d int) (int, int) {
    switch d {
    case 1:
        return 0, -1
    case 2:
        return 0, 1
    case 3:
        return -1, 0
    case 4:
        return 1, 0
    }
    return 0, 0
}

// Given a direction get the opposite direction
func GetOppositeDirection(d int) (int) {
        switch d {
        case 1:
            return 2
        case 2:
            return 1
        case 3:
            return 4
        case 4:
            return 3
        }
        log.Fatal("GetOppositeDirection got bad value", d)
        return -1
}
