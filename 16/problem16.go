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
    scanner.Scan()

    input := make([]int, 0)
    for _, r := range(scanner.Text()){
        input = append(input, int(r-'0'))
    }

    for i := 0; i < 100; i++ {
        input = Phase(input)
    }
    fmt.Println("Problem 16a:", InputString(input[:8]))
}

func Phase(inp []int) []int {
    out := make([]int, len(inp))
    for i, _ := range(inp){
        sum := 0
        for j, jv := range(inp){
            patt := PatternAt(i, j)
            sum  += jv * patt
        }
        if sum < 0 {
            sum = sum * -1
        }
        out[i] = sum % 10
    }
    return out
}

func PatternAt(out_pos, pos int) int {
    PATTERN := [...]int{0, 1, 0, -1}
    out_pos += 1
    step := (1 + float64(pos))/float64(out_pos)
    ind  := int(step) % 4
    return PATTERN[ind]
}

func InputString(i []int) string {
    s := ""
    for _, v := range(i) {
        s += strconv.Itoa(v)
    }
    return s
}
