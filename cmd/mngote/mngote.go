package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/devtin/monigote/cmd/mngote/models"
)

func main() {
	args := models.Args{}
	arg.MustParse(&args)

	fmt.Println(args.Input)
}
