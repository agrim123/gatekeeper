package main

import (
	"fmt"
	"os"

	"github.com/agrim123/gatekeeper/internal/pkg/root"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/pkg/config"
)

func main() {
	if len(os.Args) < 2 {
		panic("Invalid arg")
	}

	config.Init()

	setup.Start()

	plan := os.Args[1]

	if _, ok := root.Plans[plan]; !ok {
		panic("Invalid plan")
	}

	if len(os.Args) < 3 {
		fmt.Println(root.Plans[plan].Options)
		return
	}

	if _, ok := root.Plans[plan].Opts[os.Args[2]]; !ok {
		panic("Invalid option")
	}

	fmt.Println(root.Plans[plan].Opts[os.Args[2]].Run())
}
