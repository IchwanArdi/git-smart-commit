package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// IsGitRepo memeriksa apakah direktori aktif saat ini merupakan repositori Git.
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// HasStagedChanges memeriksa apakah ada file yang sudah di-stage (dijalankan git add)
// dan siap untuk di-commit.
func HasStagedChanges() bool {
	// "git diff --cached --quiet" mengembalikan exit code 1 jika ada file yang di-stage,
	// dan exit code 0 jika tidak ada file yang di-stage.
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err != nil {
		// Jika error tipe ExitError, kita cek exit code-nya.
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() == 1
		}
	}
	return false
}

// ExecuteCommit menjalankan "git commit -m '<pesan>'" menggunakan os/exec.
func ExecuteCommit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("gagal melakukan commit: %w", err)
	}

	return stdout.String(), nil
}
