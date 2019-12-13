package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
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
    image_raw := scanner.Text()

    image_w := 25
    image_h := 6

    n_zero_digits := -1
    problem_8a := 0


    layers := make([][][]int, 0)
    for l := 0; l * image_w * image_h < len(image_raw); l++ {
        c_zero := 0
        s_1 := 0
        s_2 := 0
        layers = append(layers, make([][]int, 0))
        for y := 0; y < image_h; y++ {
            layers[l] = append(layers[l], make([]int, 0))
            for x := 0; x < image_w; x++ {
                i := (l*image_w*image_h)+(y*image_w)+x
                v := int(image_raw[i])-48
                if v == 0 {
                    c_zero += 1
                } else if v == 1 {
                    s_1 += 1
                } else if v == 2 {
                    s_2 += 1
                }
                layers[l][y] = append(layers[l][y], v)
            }
        }
        if (c_zero < n_zero_digits || n_zero_digits == -1) {
            n_zero_digits = c_zero
            problem_8a = s_1 * s_2
        }

    }
    fmt.Println("Problem 8a:", problem_8a)

    n_layers := len(layers)
    for y := 0; y < image_h; y++ {
        for x := 0; x < image_w; x++ {
            for l := 0; l < n_layers; l++ {
                v := layers[l][y][x]
                if v != 2 {
                    if v == 0 {
                        fmt.Print(" ")
                    } else {
                        fmt.Print("X")
                    }
                    break
                }
            }
        }
        fmt.Println("")
    }
}
