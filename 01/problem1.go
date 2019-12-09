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

    parta_fuel := 0
    partb_fuel := 0
    for scanner.Scan() {
        i, err := strconv.Atoi(scanner.Text())
        if err != nil { log.Fatal(err) }
        i = (i / 3) - 2

        parta_fuel += i
        partb_fuel += i

        extra := (i / 3) - 2
        for (extra > 0) {
            partb_fuel += extra
            extra = (extra / 3) - 2
        }
    }

    fmt.Printf("Part 1a: %v\n", parta_fuel)
    fmt.Printf("Part 1b: %v\n", partb_fuel)
}
