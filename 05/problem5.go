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
    c.Input = append(c.Input, 1)
    c.Run()
    out := c.Output[len(c.Output)-1]
    fmt.Printf("Problem 5a: %v\n", out)

    c  = Computer(program)
    c.Input = append(c.Input, 5)
    c.Run()
    out  = c.Output[len(c.Output)-1]
    fmt.Printf("Problem 5b: %v\n", out)

}

type IntcodeComputer struct {
    RAM    []int
    Eip      int
    Input  []int
    Output []int
}

func Computer (code []int) (*IntcodeComputer){
    c    := &IntcodeComputer{Eip: 0}
    c.RAM = make([]int, len(code))
    copy(c.RAM, code)
    c.Input = make([]int, 0)
    c.Output = make([]int, 0)
    return c
}

func (c *IntcodeComputer) Run () {
    for {
        opcode :=  c.RAM[c.Eip] % 100
        mode_a := (c.RAM[c.Eip] / 100) % 10
        mode_b := (c.RAM[c.Eip] / 1000) % 10
        //mode_c := (c.RAM[c.Eip] / 10000) % 10

        switch opcode {
        case 1:
            a, b := 0, 0
            if (mode_a == 0){
                a = c.RAM[c.RAM[c.Eip+1]]
            } else {
                a = c.RAM[c.Eip+1]
            }

            if (mode_b == 0){
                b = c.RAM[c.RAM[c.Eip+2]]
            } else {
                b = c.RAM[c.Eip+2]
            }
            c.RAM[c.RAM[c.Eip+3]] = a + b
            c.Eip += 4
        case 2:
            a, b := 0, 0
            if (mode_a == 0){
                a = c.RAM[c.RAM[c.Eip+1]]
            } else {
                a = c.RAM[c.Eip+1]
            }

            if (mode_b == 0){
                b = c.RAM[c.RAM[c.Eip+2]]
            } else {
                b = c.RAM[c.Eip+2]
            }
            c.RAM[c.RAM[c.Eip+3]] = a * b
            c.Eip += 4
        case 3:
            // Input
            c.RAM[c.RAM[c.Eip+1]] = c.Input[0]
            c.Input = c.Input[1:]
            c.Eip += 2
        case 4:
            //Output
            c.Output = append(c.Output, c.RAM[c.RAM[c.Eip+1]])
            c.Eip += 2
        case 5:
            //Jump-if-true
            a, b := 0, 0
            if (mode_a == 0){
                a = c.RAM[c.RAM[c.Eip+1]]
            } else {
                a = c.RAM[c.Eip+1]
            }
            if (mode_b == 0){
                b = c.RAM[c.RAM[c.Eip+2]]
            } else {
                b = c.RAM[c.Eip+2]
            }

            if a != 0 {
                c.Eip = b
            } else {
                c.Eip += 3
            }
        case 6:
            //Jump-if-false
            a, b := 0, 0
            if (mode_a == 0){
                a = c.RAM[c.RAM[c.Eip+1]]
            } else {
                a = c.RAM[c.Eip+1]
            }
            if (mode_b == 0){
                b = c.RAM[c.RAM[c.Eip+2]]
            } else {
                b = c.RAM[c.Eip+2]
            }

            if a == 0 {
                c.Eip = b
            } else {
                c.Eip += 3
            }
        case 7:
            //less-than
            a, b := 0, 0
            if (mode_a == 0){
                a = c.RAM[c.RAM[c.Eip+1]]
            } else {
                a = c.RAM[c.Eip+1]
            }
            if (mode_b == 0){
                b = c.RAM[c.RAM[c.Eip+2]]
            } else {
                b = c.RAM[c.Eip+2]
            }

            if a < b {
                c.RAM[c.RAM[c.Eip+3]] = 1
            } else {
                c.RAM[c.RAM[c.Eip+3]] = 0
            }
            c.Eip += 4
        case 8:
            //less-than
            a, b := 0, 0
            if (mode_a == 0){
                a = c.RAM[c.RAM[c.Eip+1]]
            } else {
                a = c.RAM[c.Eip+1]
            }
            if (mode_b == 0){
                b = c.RAM[c.RAM[c.Eip+2]]
            } else {
                b = c.RAM[c.Eip+2]
            }

            if a == b {
                c.RAM[c.RAM[c.Eip+3]] = 1
            } else {
                c.RAM[c.RAM[c.Eip+3]] = 0
            }
            c.Eip += 4
         case 99:
            return
        default:
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
