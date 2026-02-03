package main

import (
	"context"

	"github.com/DimaKropachev/cryptool/internal/cli"
)

func main() {
	ctx := context.Background()

	cli.Execute(ctx)
}
