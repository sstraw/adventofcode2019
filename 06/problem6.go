package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
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

    orbits := make(map[string][]string)

    for scanner.Scan() {
        t := strings.Split(scanner.Text(), ")")
        orbits[t[0]] = append(orbits[t[0]], t[1])
    }

    sum := 0
    for k, _ := range(orbits) {
        sum += (CalculateOrbits(orbits, k)-1)
    }

    fmt.Printf("Problem 6a: %v\n", sum)
}

func CalculateOrbits(orbits map[string][]string, orbit string) int {
    s, exists := orbits[orbit]
    if (!exists) {
        return 1
    }
    t := 1
    for _, o := range(s) {
        t += CalculateOrbits(orbits, o)
    }
    return t
}
