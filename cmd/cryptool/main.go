/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/DimaKropachev/cryptool/internal/cli"
	"github.com/DimaKropachev/cryptool/pkg/logger"
)

func main() {
	ctx, err := logger.New(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialization logger")
		os.Exit(0)
	}

	cli.Execute(ctx)
}
