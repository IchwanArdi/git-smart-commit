package cmd

import (
	"fmt"
	"os"

	"github.com/IchwanArdi/git-smart-commit/internal/git"
	"github.com/IchwanArdi/git-smart-commit/internal/prompt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-smart-commit",
	Short: "git-smart-commit adalah CLI tool untuk membuat Conventional Commit secara interaktif",
	Long: `Sebuah aplikasi CLI pembantu untuk membuat commit message yang rapi, 
konsisten, dan mengikuti standar Conventional Commits secara interaktif.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Validasi apakah ini folder repositori Git
		if !git.IsGitRepo() {
			fmt.Println("⚠️  Error: Folder ini bukan repositori Git. Silakan jalankan 'git init' terlebih dahulu.")
			os.Exit(1)
		}

		// 2. Validasi apakah ada perubahan yang siap di-commit (staged changes)
		if !git.HasStagedChanges() {
			fmt.Println("⚠️  Peringatan: Tidak ada perubahan yang telah di-stage (staged).")
			fmt.Println("Silakan jalankan 'git add <file>' terlebih dahulu sebelum membuat commit.")
			os.Exit(0)
		}

		// 3. Ambil file yang di-stage & tebak scope-nya secara otomatis
		stagedFiles, err := git.GetStagedFiles()
		guessedScope := ""
		if err == nil {
			guessedScope = git.GuessScope(stagedFiles)
		}

		// 4. Tampilkan TUI Form (Huh)
		answers, err := prompt.AskQuestions(guessedScope)
		if err != nil {
			fmt.Printf("❌ Batal: %v\n", err)
			os.Exit(1)
		}

		// 5. Cek konfirmasi commit dari user
		if !answers.ConfirmCommit {
			fmt.Println("🚫 Commit dibatalkan.")
			os.Exit(0)
		}

		// 6. Format dan eksekusi commit
		commitMsg := answers.FormatCommitMessage()
		output, err := git.ExecuteCommit(commitMsg)
		if err != nil {
			fmt.Printf("❌ Error saat commit: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Berhasil membuat commit!")
		fmt.Println(output)
	},
}

// Execute menjalankan CLI tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
