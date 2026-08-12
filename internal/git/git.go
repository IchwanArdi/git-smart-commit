package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() == 1
		}
	}
	return false
}

// GetStagedFiles mengambil daftar path file yang sedang di-stage saat ini.
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	var files []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}
	return files, nil
}

// GuessScope menebak scope terbaik berdasarkan nama direktori file yang baru saja di-stage.
func GuessScope(files []string) string {
	if len(files) == 0 {
		return ""
	}

	// Ambil file pertama sebagai referensi
	firstFile := files[0]

	// Pisahkan path file berdasarkan slash "/" (atau backslash "\" di Windows jika diformat git)
	// Git path selalu menggunakan "/" forward slash
	parts := strings.Split(firstFile, "/")
	if len(parts) > 1 {
		// Mengambil nama folder terdekat sebelum nama file
		// Contoh: internal/git/git.go -> git
		return parts[len(parts)-2]
	}

	// Jika file berada di root (misal main.go), kita kosongkan default scope-nya
	return ""
}

// ExecuteCommit menjalankan "git commit -m '<pesan>'".
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

// GetRemotes mengambil daftar remote yang dikonfigurasi.
func GetRemotes() ([]string, error) {
	cmd := exec.Command("git", "remote")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	var remotes []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			remotes = append(remotes, trimmed)
		}
	}
	return remotes, nil
}

// HasRemote memeriksa apakah repositori memiliki remote yang dikonfigurasi.
func HasRemote() bool {
	remotes, err := GetRemotes()
	return err == nil && len(remotes) > 0
}

// GetCurrentBranch mengambil nama branch saat ini.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// ExecutePush menjalankan "git push <remote> <branch>".
func ExecutePush(remote, branch string) (string, error) {
	cmd := exec.Command("git", "push", remote, branch)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String() + stdout.String(), fmt.Errorf("gagal melakukan push: %w", err)
	}

	return stdout.String() + stderr.String(), nil
}

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	return string(output), nil
}
