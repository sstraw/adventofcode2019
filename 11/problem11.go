package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strconv"
    "bytes"
    "strings"
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
    program := make([]int, 0)
    for _, i := range(strings.Split(scanner.Text(),",")){
        i, _ := strconv.Atoi(i)
        program = append(program, i)
    }

    r := NewPaintRobot(program)
    r.Run()
    n_squares := 0
    for _, ym := range(r.H.P) {
        n_squares += len(ym)
    }
    fmt.Println("Problem 11a:", n_squares)

    r  = NewPaintRobot(program)
    r.H.WriteSquare(0, 0, 1)
    r.Run()
    min_x, max_x, min_y, max_y := 0, 0, 0, 0
    for x, xm := range(r.H.P){
        if x < min_x {
            min_x = x
        }
        if x > max_x {
            max_x = x
        }
        for y, _ := range(xm){
            if y < min_y {
                min_y = y
            }
            if y > max_y {
                max_y = y
            }
        }
    }

    fmt.Println("Problem 11b:")
    for y := max_y; y >= min_y; y -= 1 {
        for x := min_x; x <= max_x; x++ {
            switch r.H.ReadSquare(x, y){
            case 0:
                fmt.Print(" ")
            case 1:
                fmt.Print("X")
            default:
                log.Fatal("Shouldn't be here")
            }
        }
        fmt.Println()
    }
}

type PaintRobot struct {
    X int
    Y int
    D int // Direction
    C *IntcodeComputer
    H *Hull
}

func NewPaintRobot (program []int) (*PaintRobot) {
    return &PaintRobot {
        X: 0,
        Y: 0,
        D: 0,
        C: Computer(program),
        H: NewHull(),
    }
}

func (pr *PaintRobot) Run() {
    go pr.C.Run()

    for {
        val := pr.H.ReadSquare(pr.X, pr.Y)
        pr.C.Input <- val
        new_val, more := <-pr.C.Output
        if !more {
            break
        }
        new_dir := <-pr.C.Output

        pr.H.WriteSquare(pr.X, pr.Y, new_val)

        if new_dir == 1 {
            pr.D += 1
        } else {
            pr.D -= 1
        }
        if pr.D == -1 {
            pr.D = 3
        } else {
            pr.D = pr.D % 4
        }

        switch pr.D {
        case 0:
            pr.Y += 1
        case 1:
            pr.X += 1
        case 2:
            pr.Y -= 1
        case 3:
            pr.X -= 1
        default:
            log.Fatal("Direction value: ", pr.D)
        }
    }
}

type Hull struct {
    P map[int]map[int]int
}

func NewHull () (*Hull) {
    return &Hull{
        make(map[int]map[int]int),
    }
}

func (hull *Hull) ReadSquare(x, y int) int {
    if xm, ok := hull.P[x]; ok {
        if v, ok := xm[y]; ok {
            return v
        } else {
            return 0
        }
    } else {
        return 0
    }
}

func (hull *Hull) WriteSquare(x, y, v int) {
    if _, ok := hull.P[x]; !ok {
        hull.P[x] = make(map[int]int)
    }
    hull.P[x][y] = v
}


type IntcodeComputer struct {
    ROM    []int
    RAM    []int
    Eip      int
    Input    chan int
    Output   chan int
    ExitCode int
    Running  bool
    RelBase  int
}

// A struct to represent a basic opcode
// A, B, and C should always represent the location
// in memory where the appropriate data can be found.
type Instruction struct {
    Opcode int
    A      int
    B      int
    C      int
}

func (c *IntcodeComputer) GetOperation(eip int) (*Instruction) {
    opcode :=  c.RAM[eip] % 100
    mode   := make([]int, 3)
    mode[0] = (c.RAM[eip] / 100) % 10
    mode[1] = (c.RAM[eip] / 1000) % 10
    mode[2] = (c.RAM[eip] / 10000) % 10

    vals   := make([]int, 3)

    for i:= 0; i < 3; i++{
        // Handle if blind arg processing takes us past the end
        // It is up to the instruction section to deal with this
        if (eip+1+i) >= len(c.RAM) {
            vals[i] = -1
            continue
        }
        switch mode[i] {
        case 0:
            // Position mode
            vals[i] = c.RAM[eip+1+i]
        case 1:
            // Immediate mode
            vals[i] = eip+1+i
        case 2:
            // Relative mode
            vals[i] = c.RelBase + c.RAM[eip+1+i]
        }
    }

    inst := & Instruction{opcode, vals[0], vals[1], vals[2]}
    return inst
}

func Computer (code []int) (*IntcodeComputer){
    c    := &IntcodeComputer{Eip: 0, Running: false, RelBase: 0}
    c.ROM = make([]int, len(code))
    copy(c.ROM, code)
    c.Input  = make(chan int, 1)
    c.Output = make(chan int, 2)
    c.RAM = make([]int, len(code)*999)
    copy(c.RAM, code)
    return c
}

func (c *IntcodeComputer) Run () {
    copy(c.RAM, c.ROM)
    c.Running = true
    for {
        inst := c.GetOperation(c.Eip)

        switch inst.Opcode {
        case 1:
            c.RAM[inst.C] = c.RAM[inst.A] + c.RAM[inst.B]
            c.Eip += 4
        case 2:
            c.RAM[inst.C] = c.RAM[inst.A] * c.RAM[inst.B]
            c.Eip += 4
         case 3:
            // Input
            c.RAM[inst.A] = <-c.Input
            c.Eip += 2
        case 4:
            //Output
            c.Output <- c.RAM[inst.A]
            if c.RAM[c.Eip+2] == 99 {
                c.ExitCode = c.RAM[inst.A]
            }
            c.Eip += 2
        case 5:
            //Jump-if-true
            if c.RAM[inst.A] != 0 {
                c.Eip = c.RAM[inst.B]
            } else {
                c.Eip += 3
            }
        case 6:
            //Jump-if-false
            if c.RAM[inst.A] == 0 {
                c.Eip = c.RAM[inst.B]
            } else {
                c.Eip += 3
            }
        case 7:
            //less-than
            if c.RAM[inst.A] < c.RAM[inst.B] {
                c.RAM[inst.C] = 1
            } else {
                c.RAM[inst.C] = 0
            }
            c.Eip += 4
        case 8:
            //equals
            if c.RAM[inst.A] == c.RAM[inst.B] {
                c.RAM[inst.C] = 1
            } else {
                c.RAM[inst.C] = 0
            }
            c.Eip += 4
        case 9:
            // Relative base adjustment
            c.RelBase += c.RAM[inst.A]
            c.Eip += 2
        case 99:
            c.Running = false
            close(c.Output)
            return
        default:
            c.Running = false
            close(c.Output)
            return
        }
    }
}

func ScanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
  if atEOF && len(data) == 0 {
    return 0, nil, nil
  }
  if i := bytes.IndexByte(data, ','); i >= 0 {
    // We have a full newline-terminated line.
    return i + 1, data[0:i], nil
  }
  // If we're at EOF, we have a final, non-terminated line. Return it.
  if atEOF {
    return len(data), data, nil
  }
  // Request more data.
  return 0, nil, nil
}
