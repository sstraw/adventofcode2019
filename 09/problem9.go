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

    c := Computer(program)
    c.Input <- 1
    c.Run()
    fmt.Println("Problem 9a:", c.ExitCode)
}

// Perm calls f with each permutation of a.
func Perm(a []int, f func([]int)) {
    perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, f func([]int), i int) {
    if i > len(a) {
        f(a)
        return
    }
    perm(a, f, i+1)
    for j := i + 1; j < len(a); j++ {
        a[i], a[j] = a[j], a[i]
        perm(a, f, i+1)
        a[i], a[j] = a[j], a[i]
    }
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
    c.Input  = make(chan int, 5)
    c.Output = make(chan int, 5)
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
