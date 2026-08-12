package main

import (
	"fmt"
	"log"

	"github.com/IchwanArdi/git-smart-commit/internal/ai"
	gitservice "github.com/IchwanArdi/git-smart-commit/internal/git"
)

func main() {
	diff, err := gitservice.GetStagedDiff()
	if err != nil {
		log.Fatal(err)
	}

	if diff == "" {
		fmt.Println("Tidak ada staged changes.")
		return
	}

	fmt.Println("🔍 Menganalisis staged changes...")
	fmt.Println()

	message, err := ai.GenerateCommitMessage(diff)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🤖 Suggested commit:")
	fmt.Println()
	fmt.Println(message)
}
