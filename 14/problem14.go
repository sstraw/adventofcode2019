package main

import (
    "fmt"
    "os"
    "bufio"
    "log"
    "strconv"
    "regexp"
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

    ing_rex := regexp.MustCompile(`\d+ \w+`)

    recipes := make([]*Recipe, 0)
    for scanner.Scan() {
        recipe := &Recipe{}
        recipe.Components = make(map[string]int, 0)
        r_rec  := strings.Split(scanner.Text(), "=")
        res_m  := ing_rex.FindString(r_rec[1])
        res, res_n := ParseIngredient(res_m)
        recipe.Amount = res_n
        r_com  := strings.Split(r_rec[0], ",")
        for _, ing := range(r_com) {
            ing_s, ing_n := ParseIngredient(ing)
            recipe.Components[ing_s] = ing_n
        }
        recipe.Name = res
        recipes = append(recipes, recipe)
    }

    chemicals       := make(map[string]bool, 0)
    ord_recipes     := make([]*Recipe, 0)
    chemicals["ORE"] = true


    // Sort based on dependencies 
    for len(ord_recipes) < len(recipes) {
        for _, r := range(recipes) {
            if _, ok := chemicals[r.Name]; !ok {
                can_add := true
                for c, _ := range(r.Components) {
                    if _, ok := chemicals[c]; !ok {
                        can_add = false
                        break
                    }
                }
                if can_add {
                    ord_recipes = append(ord_recipes, r)
                    chemicals[r.Name] = true
                }
            }
        }
    }

    // 14a
    chem_need := make(map[string]int, 0)
    chem_need["FUEL"] = 1
    for i := len(ord_recipes) - 1; i >= 0; i-- {
        r    := ord_recipes[i]
        mult := chem_need[r.Name] / r.Amount
        if chem_need[r.Name] % r.Amount != 0 {
            mult += 1
        }
        for cs, cn := range(r.Components) {
            if _, ok := chem_need[cs]; !ok {
                chem_need[cs] = 0
            }
            chem_need[cs] += cn * mult
        }
    }
    fmt.Println("Problem 14a:", chem_need["ORE"])
}

type Recipe struct {
    Amount int
    Name string
    Components map[string]int
}

type Nanofactory struct {
    Requirements map[string]int
    Chemicals []string //List of all components sorted by dependencies
    Recipes map[string]*Recipe
}

func NewNanofactory () *Nanofactory {
    return &Nanofactory {
        Requirements: make(map[string]int, 0),
        Chemicals: make([]string, 0),
        Recipes: make(map[string]*Recipe, 0),
    }
}

func ParseIngredient(s string) (string, int) {
    s   = strings.Trim(s, " ")
    sp := strings.Split(s, " ")
    n, _ := strconv.Atoi(sp[0])
    return sp[1], n
}
