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

    matches := make([]int, 0)

    for i := low; i < high; i+=1 {
        // Turn it into a list
        j := make([]int, 6)
        j[0] =  i / 100000
        j[1] = (i / 10000) % 10
        j[2] = (i / 1000) % 10
        j[3] = (i / 100) % 10
        j[4] = (i / 10) % 10
        j[5] = (i / 1) % 10

        // Check if it matches criteria
        if ((j[0] == j[1] ||
             j[1] == j[2] ||
             j[2] == j[3] ||
             j[3] == j[4] ||
             j[4] == j[5]) &&
            (j[0] <= j[1] &&
             j[1] <= j[2] &&
             j[2] <= j[3] &&
             j[3] <= j[4] &&
             j[4] <= j[5])){
            matches = append(matches, i)
        }
        //Generate guess based on incrementing digits
        
    }
    fmt.Printf("Problem 4a: %v\n", len(matches))
}
