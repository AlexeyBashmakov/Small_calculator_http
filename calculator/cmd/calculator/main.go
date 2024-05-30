package main

import (
	"context"
	"fmt"
	"os"

	"calculator/internal/application"
	"calculator/internal/environ_vars"
)

func greeting() {
	prnt("Привет, апельсин :)")
}

func prnt(a ...any) {
	fmt.Println(a...)
}

func main() {
	greeting()

	ctx := context.Background()
	// Exit приводит к завершению программы с заданным кодом.
	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) int {
	if !environ_vars.CheckEnvironmentVariables() {
		if environ_vars.SetEnvironmentVariables() != nil {
			return 2
		}
	}
	return application.Run(ctx)
}
