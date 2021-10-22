package main

import (
	"github.com/igridnet/igrid/cli"
	"log"
	"os"
)

func main() {
	app := cli.New()

	if err := app.Run(os.Args); err != nil{
		log.Fatalf("%v\n", err)
	}
}
