package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	arch, err := getArch()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	depsInstalled, err := depsInstalled()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(arch)
	fmt.Println(depsInstalled)
}
