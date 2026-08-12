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

		// 2. Validasi staged changes
		if !git.HasStagedChanges() {
			fmt.Println("⚠️  Peringatan: Tidak ada perubahan yang telah di-stage (staged).")
			fmt.Println("Silakan jalankan 'git add <file>' terlebih dahulu sebelum membuat commit.")
			os.Exit(0)
		}

		// 3. Ambil staged files & tebak scope
		stagedFiles, err := git.GetStagedFiles()
		guessedScope := ""

		if err == nil {
			guessedScope = git.GuessScope(stagedFiles)
		}

		// 4. Ambil branch
		branchName := ""

		if git.HasRemote() {
			if b, err := git.GetCurrentBranch(); err == nil {
				branchName = b
			}
		}

		// =========================================================
		// 5. Generate commit message menggunakan AI
		// =========================================================

		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Printf("❌ Gagal mengambil staged diff: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🤖 Menganalisis perubahan dengan Gemini...")

		aiCommitMsg, err := ai.GenerateCommitMessage(diff)

		var commitMsg string
		var confirmPush bool

		if err == nil {
			// AI berhasil
			fmt.Println()
			fmt.Println("🤖 Suggested commit:")
			fmt.Println()
			fmt.Println(aiCommitMsg)
			fmt.Println()

			useAI, err := prompt.AskAICommitConfirmation(aiCommitMsg)
			if err != nil {
				fmt.Printf("❌ Batal: %v\n", err)
				os.Exit(1)
			}

			if useAI {
				// Gunakan hasil AI
				commitMsg = strings.TrimSpace(aiCommitMsg)

				// Tanya push secara terpisah
				confirmPush, err = prompt.AskPushConfirmation(branchName)
				if err != nil {
					fmt.Printf("❌ Batal: %v\n", err)
					os.Exit(1)
				}

			} else {
				// User menolak hasil AI → masuk mode manual
				fmt.Println()
				fmt.Println("✍️  Beralih ke mode manual...")

				answers, err := prompt.AskQuestions(guessedScope, branchName)
				if err != nil {
					fmt.Printf("❌ Batal: %v\n", err)
					os.Exit(1)
				}

				if !answers.ConfirmCommit {
					fmt.Println("🚫 Commit dibatalkan.")
					os.Exit(0)
				}

				commitMsg = answers.FormatCommitMessage()
				confirmPush = answers.ConfirmPush
			}

		} else {
			// Gemini gagal → fallback ke mode manual
			fmt.Printf("⚠️  Gemini tidak dapat digunakan: %v\n", err)
			fmt.Println("✍️  Beralih ke mode manual...")

			answers, err := prompt.AskQuestions(guessedScope, branchName)
			if err != nil {
				fmt.Printf("❌ Batal: %v\n", err)
				os.Exit(1)
			}

			if !answers.ConfirmCommit {
				fmt.Println("🚫 Commit dibatalkan.")
				os.Exit(0)
			}

			commitMsg = answers.FormatCommitMessage()
			confirmPush = answers.ConfirmPush
		}

		// =========================================================
		// 6. Execute commit
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
		// 7. Push jika dikonfirmasi
		// =========================================================

		if confirmPush {
			fmt.Println("🚀 Sedang melakukan push ke remote repository...")

			// Pastikan branchName tidak kosong
			if branchName == "" {
				branchName, err = git.GetCurrentBranch()
				if err != nil {
					fmt.Printf("❌ Gagal mendapatkan nama branch saat ini: %v\n", err)
					os.Exit(1)
				}
			}

			// Ambil remote
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

			fmt.Printf("✅ Berhasil push ke %s/%s!\n", remote, branchName)

			if strings.TrimSpace(pushOutput) != "" {
				fmt.Println(pushOutput)
			}
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
