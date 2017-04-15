package main

import (
	"github.com/DanShu93/onion-of-life/game"
	"flag"
	"strings"
	"strconv"
)

func main() {
	var rulesString string
	flag.StringVar(&rulesString, "rule", "23/3", "The rules to apply e.g. 23/3")

	flag.Parse()

	args := flag.Args()

	rules := strings.Split(rulesString, "/")

	aliveAmounts := toInts(strings.Split(rules[0], ""))
	bornAmounts := toInts(strings.Split(rules[1], ""))

	controller := game.NewController(args[0], aliveAmounts, bornAmounts)
	controller.Play()
}

func toInts(strings []string) []int {
	ints := make([]int, len(strings))

	for i, s := range strings {
		number, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}

		ints[i] = number
	}

	return ints
}
