package cmd

import (
	"fmt"
	"os"

	"git-smart-commit/internal/git"
	"git-smart-commit/internal/prompt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-smart",
	Short: "git-smart adalah CLI tool untuk membuat Conventional Commit secara interaktif",
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

		// 3. Tampilkan form pertanyaan interaktif
		answers, err := prompt.AskQuestions()
		if err != nil {
			fmt.Printf("❌ Batal: %v\n", err)
			os.Exit(1)
		}

		// 4. Format pesan commit
		commitMsg := answers.FormatCommitMessage()

		// 5. Tampilkan preview commit ke pengguna
		fmt.Println("\n📝 Preview Pesan Commit:")
		fmt.Println("========================================")
		fmt.Println(commitMsg)
		fmt.Println("========================================")

		// 6. Konfirmasi sebelum commit dilakukan
		confirm := false
		confirmPrompt := &survey.Confirm{
			Message: "Apakah Anda ingin membuat commit dengan pesan di atas?",
			Default: true,
		}
		if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
			fmt.Printf("❌ Batal: %v\n", err)
			os.Exit(1)
		}

		if !confirm {
			fmt.Println("🚫 Commit dibatalkan.")
			os.Exit(0)
		}

		// 7. Eksekusi commit
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
