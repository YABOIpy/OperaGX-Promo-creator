package main

import (
	"fmt"
	"promogen/internal/generator"
	"promogen/internal/utils"
	"time"
)

func main() {
	instances, err := generator.CreateInstances()
	if err != nil {
		panic(err)
	}
	var tokens []string
	start := time.Now()
	generator.StartTask(instances, func(c generator.Instance) {
		s := time.Now()
		if token := c.GetOperaToken(); token != "" {
			fmt.Printf(utils.LogFormat, time.Since(s).Milliseconds(), token[:40]+"****")
			tokens = append(tokens, generator.PromoURL+token)
		}
	})
	fmt.Printf(utils.CountFormat, time.Since(start).Seconds(), len(tokens))

	utils.WriteArrayToFile("promos.txt", tokens)
}
