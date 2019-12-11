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

    max_thrust := 0
    settings := []int{0,1,2,3,4}
    Perm(settings, func(a []int){
        thrust := 0
        for i := 0; i < 5; i++ {
            c := Computer(program)
            c.Input <- settings[i]
            c.Input <- thrust
            go c.Run()
            thrust = <-c.Output
        }
        if thrust > max_thrust {
            max_thrust = thrust
        }
    })

    fmt.Println("Problem 7a:", max_thrust)

    max_thrust = 0
    settings   = []int{5,6,7,8,9}
    Perm(settings, func(a []int){
        thrusters := make([]*IntcodeComputer,0)
        for i:=0; i<5; i++ {
            thrusters = append(thrusters, Computer(program))
            thrusters[i].Input <- settings[i]
        }
        for i:=0; i<4; i++ {
            thrusters[i].Output = thrusters[i+1].Input
        }
        for i := 0; i < 5; i++ {
            go thrusters[i].Run()
        }
        thrusters[0].Input <- 0

        x := 0
        for {
            x = <-thrusters[4].Output
            if thrusters[4].Running {
                thrusters[0].Input <- x
            } else{
                break
            }
        }
        if x > max_thrust { max_thrust = x }
    })
    fmt.Println("Problem 7b:", max_thrust)
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
    RAM    []int
    Eip      int
    Input    chan int
    Output   chan int
    ExitCode int
    Running  bool
}

func Computer (code []int) (*IntcodeComputer){
    c    := &IntcodeComputer{Eip: 0, Running: false}
    c.RAM = make([]int, len(code))
    copy(c.RAM, code)
    c.Input  = make(chan int, 5)
    c.Output = make(chan int, 5)
    return c
}

func (c *IntcodeComputer) Run () {
    c.Running = true
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
            c.RAM[c.RAM[c.Eip+1]] = <-c.Input
            c.Eip += 2
        case 4:
            //Output
            c.Output <- c.RAM[c.RAM[c.Eip+1]]
            if c.RAM[c.Eip+2] == 99 {
                c.ExitCode = c.RAM[c.RAM[c.Eip+1]]
            }
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
            c.Running = false
            return
        default:
            c.Running = false
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
