package main

import (
	_ "github.com/joho/godotenv/autoload"
	_ "go.uber.org/automaxprocs"

	"github.com/dingdayu/dnsx/cmd"
)

func main() {
	cmd.Execute()
}
