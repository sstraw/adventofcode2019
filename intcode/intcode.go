package intcode

import (
    "fmt"
    "strconv"
    "strings"
)

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

func (i Instruction) String() string {
    var s_opcode string
    switch i.Opcode {
    case 1:
        s_opcode = "add"
    case 2:
        s_opcode = "mult"
    case 3:
        s_opcode = "input"
    case 4:
        s_opcode = "output"
    case 5:
        s_opcode = "jt"
    case 6:
        s_opcode = "jf"
    case 7:
        s_opcode = "lt"
    case 8:
        s_opcode = "eq"
    case 9:
        s_opcode = "adjbase"
    }
    return fmt.Sprintf("%v %v %v %v", s_opcode, i.A, i.B, i.C)
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
    //fmt.Println(c.Eip, inst)
    return inst
}

func NewComputer (code string) (*IntcodeComputer){
    c    := &IntcodeComputer{Eip: 0, Running: false, RelBase: 0}
    c.ROM = make([]int, 0)

    for _, i := range(strings.Split(code,",")){
        i, _ := strconv.Atoi(i)
        c.ROM = append(c.ROM, i)
    }

    c.Input  = make(chan int)
    c.Output = make(chan int)
    c.RAM = make([]int, len(c.ROM)*999)
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
