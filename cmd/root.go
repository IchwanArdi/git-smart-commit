package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/IchwanArdi/git-smart-commit/internal/ai"
	"github.com/IchwanArdi/git-smart-commit/internal/git"
	"github.com/IchwanArdi/git-smart-commit/internal/prompt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-smart",
	Short: "git-smart-commit adalah AI-powered CLI tool untuk membuat Conventional Commit",
	Long: `Git-Smart-Commit adalah CLI tool yang menggunakan Gemini AI
untuk menganalisis staged Git changes dan menghasilkan commit message
berdasarkan standar Conventional Commits.`,

	Run: func(cmd *cobra.Command, args []string) {

		// 1. Validasi apakah folder saat ini merupakan repositori Git
		if !git.IsGitRepo() {
			fmt.Println("❌ Error: Folder ini bukan repositori Git.")
			fmt.Println("Silakan jalankan 'git init' terlebih dahulu.")
			os.Exit(1)
		}

		// 2. Validasi apakah terdapat staged changes
		if !git.HasStagedChanges() {
			fmt.Println("⚠️  Tidak ada perubahan yang telah di-stage.")
			fmt.Println("Silakan jalankan 'git add <file>' terlebih dahulu.")
			os.Exit(0)
		}

		// 3. Ambil staged diff
		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Printf("❌ Gagal mengambil staged diff: %v\n", err)
			os.Exit(1)
		}

		// 4. Ambil branch saat ini
		branchName := ""

		if git.HasRemote() {
			if b, err := git.GetCurrentBranch(); err == nil {
				branchName = b
			}
		}

		// =========================================================
		// 5. Generate commit message menggunakan Gemini
		// =========================================================

		fmt.Println("🤖 Menganalisis perubahan dengan Gemini...")

		suggestion, err := ai.GenerateCommitSuggestion(diff)
		if err != nil {
			fmt.Printf("\n❌ Gagal menghasilkan commit message:\n%v\n", err)
			os.Exit(1)
		}

		// Format structured AI response menjadi Conventional Commit
		commitMsg := strings.TrimSpace(
			ai.FormatCommitMessage(suggestion),
		)

		// =========================================================
		// 6. Tampilkan hasil AI
		// =========================================================

		fmt.Println()
		fmt.Println("🤖 Suggested commit:")
		fmt.Println()
		fmt.Println(commitMsg)
		fmt.Println()

		// =========================================================
		// 7. Konfirmasi commit
		// =========================================================

		useAI, err := prompt.AskAICommitConfirmation()
		if err != nil {
			fmt.Printf("❌ Batal: %v\n", err)
			os.Exit(1)
		}

		if !useAI {
			fmt.Println("🚫 Commit dibatalkan.")
			os.Exit(0)
		}

		// =========================================================
		// 8. Execute commit
		// =========================================================

		output, err := git.ExecuteCommit(commitMsg)
		if err != nil {
			fmt.Printf("❌ Error saat commit: %v\n", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println("✅ Berhasil membuat commit!")
		fmt.Println(output)

		// =========================================================
		// 9. Konfirmasi push
		// =========================================================

		confirmPush, err := prompt.AskPushConfirmation(branchName)
		if err != nil {
			fmt.Printf("❌ Batal: %v\n", err)
			os.Exit(1)
		}

		if !confirmPush {
			return
		}

		// =========================================================
		// 10. Execute push
		// =========================================================

		fmt.Println("🚀 Sedang melakukan push ke remote repository...")

		// Pastikan branchName tersedia
		if branchName == "" {
			branchName, err = git.GetCurrentBranch()
			if err != nil {
				fmt.Printf(
					"❌ Gagal mendapatkan nama branch saat ini: %v\n",
					err,
				)
				os.Exit(1)
			}
		}

		// Ambil remote repository
		remotes, err := git.GetRemotes()
		remote := "origin"

		if err == nil && len(remotes) > 0 {
			remote = remotes[0]
		}

		pushOutput, err := git.ExecutePush(remote, branchName)
		if err != nil {
			fmt.Printf(
				"❌ Gagal melakukan push ke %s/%s:\n%s\n",
				remote,
				branchName,
				err,
			)
			os.Exit(1)
		}

		fmt.Printf(
			"✅ Berhasil push ke %s/%s!\n",
			remote,
			branchName,
		)

		if strings.TrimSpace(pushOutput) != "" {
			fmt.Println(pushOutput)
		}
	},
}

// Execute menjalankan CLI tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
