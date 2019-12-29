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

    output := make([]int, len(input))
    copy(output, input)
    for i := 0; i < 100; i++ {
        output = Phase(output)
    }
    fmt.Println("Problem 16a:", InputString(output[:8]))

    msg_offset, _ := strconv.Atoi(scanner.Text()[:7])
    large_input   := make([]int, len(input)*10000)
    for i := 0; i < 10000*len(input); i += len(input) {
        copy(large_input[i:], input)
    }


    // As observed in the subreddit - for all indices 
    // where i > len/2, the pattern will just be 
    // 1, 1, 1, 1, etc

    if (msg_offset <= len(large_input)/2 ||
        msg_offset >= len(large_input)) {
        log.Fatal("Message offset doesn't match expected")
    }
    for n := 0; n < 100; n++ {
        sum := 0
        for i := len(large_input) - 1; i >= msg_offset; i-- {
            //fmt.Println("i: ",i)
            sum += large_input[i]
            large_input[i] = sum % 10
            if large_input[i] < 0 {
                large_input[i] = large_input[i] * -1
            }
        }
    }
    str := InputString(large_input[msg_offset:msg_offset+8])
    fmt.Println("Problem 16b:", str)
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
