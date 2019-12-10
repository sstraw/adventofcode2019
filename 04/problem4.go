package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strconv"
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
    s := scanner.Text()

    s_s := strings.Split(s, "-")

    low,  _ := strconv.Atoi(s_s[0])
    high, _ := strconv.Atoi(s_s[1])

    fmt.Printf("Range: %v - %v \n", low, high)

    matches_a := make([]int, 0)
    matches_b := make([]int, 0)

    for i := low; i < high; i+=1 {
        // Turn it into a list
        j := make([]int, 6)
        j[0] =  i / 100000
        j[1] = (i / 10000) % 10
        j[2] = (i / 1000) % 10
        j[3] = (i / 100) % 10
        j[4] = (i / 10) % 10
        j[5] = (i / 1) % 10

        // Criteria
        var n [10]int
        for _, x := range(j) {
            n[x] += 1
        }
        two_only := false
        two_or_more := false
        for _, x := range(n) {
            if (x == 2) {
                two_only = true
            }
            if (x <= 2) {
                two_or_more = true
            }
        }
        if ((j[0] <= j[1] &&
             j[1] <= j[2] &&
             j[2] <= j[3] &&
             j[3] <= j[4] &&
             j[4] <= j[5])){
            if (two_only) {
                matches_a = append(matches_a, i)
                matches_b = append(matches_b, i)
            } else if (two_or_more) {
                matches_a = append(matches_a, i)
            }
        }
    }
    fmt.Printf("Problem 4a: %v\n", len(matches_a))
    fmt.Printf("Problem 4b: %v\n", len(matches_b))
}
