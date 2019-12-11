package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strings"
)

type orbit struct {
    Planet string
    Center string
    Orbits []string
}

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

    orbits := make(map[string]orbit)

    san_step := ""
    you_step := ""
    for scanner.Scan() {
        // x)y : x orbited by y
        t := strings.Split(scanner.Text(), ")")
        t_0, t_0_ok := orbits[t[0]]
        t_1, t_1_ok := orbits[t[1]]

        if !t_0_ok {
            t_0 = orbit{t[0], "", nil}
        }
        t_0.Orbits = append(t_0.Orbits, t[1])
        orbits[t[0]] = t_0

        if !t_1_ok {
            t_1 = orbit{t[1], t[0], nil}
        }
        if (t_1.Center == "") {
            t_1.Center = t[0]
        }
        orbits[t[1]] = t_1

        if (t[1] == "SAN") {
            san_step = t[0]
        } else if (t[1] == "YOU") {
            you_step = t[0]
        }
    }

    sum := 0
    for k, _ := range(orbits) {
        sum += (CalculateOrbits(orbits, k)-1)
    }

    fmt.Printf("Problem 6a: %v\n", sum)

    san_path := []string{san_step}
    for san_step != "" {
        o := orbits[san_step]
        san_path = append(san_path, o.Center)
        san_step = o.Center
    }

    you_path := []string{you_step}
    for you_step != "" {
        o := orbits[you_step]
        you_path = append(you_path, o.Center)
        you_step = o.Center
    }

    tot := 0
    for si, sv := range(you_path) {
        for yi, yv := range(san_path){
            if sv == yv {
                tot = si + yi
                break
            }
        }
        if tot != 0 {
            break
        }
    }
    fmt.Printf("Problem 6b: %v\n", tot)
}

func CalculateOrbits(orbits map[string]orbit, orb string) int {
    s, exists := orbits[orb]
    if (!exists) {
        return 1
    }
    t := 1
    for _, o := range(s.Orbits) {
        t += CalculateOrbits(orbits, o)
    }
    return t
}
