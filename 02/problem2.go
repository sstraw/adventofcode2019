package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strconv"
    "bytes"
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
    scanner.Split(ScanCommas)

    program := make([]int, 0)
    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        program = append(program, i)
    }

    // hardcoded as per instructions, if real input
    if (fname == "input.txt"){
        program[1] = 12
        program[2] = 2
    }

    // Instruction we're on
    eip := 0
    for {
        opcode := program[eip]
        if (opcode==1) {
            program[program[eip+3]] = program[program[eip+1]] + program[program[eip+2]]
        } else if (opcode == 2) {
            program[program[eip+3]] = program[program[eip+1]] * program[program[eip+2]]
        } else if (opcode == 99) {
            break
        } else {
            break
        }
        eip += 4
    }

    fmt.Printf("Problem 2a: %v\n", program[0])

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
