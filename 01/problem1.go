package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strconv"
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

    sum := 0
    for scanner.Scan() {
        i, err := strconv.Atoi(scanner.Text())
        if err != nil { log.Fatal(err) }
        i = (i / 3) - 2
        sum += i
    }
    fmt.Printf("Sum: %v\n", sum)
}
