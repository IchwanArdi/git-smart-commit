package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// CommitAnswers menyimpan jawaban dari user untuk setiap bagian Conventional Commits.
type CommitAnswers struct {
	Type          string
	Scope         string
	Subject       string
	Body          string
	ConfirmCommit bool
	ConfirmPush   bool
}

// AskQuestions memicu formulir interaktif satu halaman (TUI Form) menggunakan Huh.
func AskQuestions(defaultScope string, branchName string) (*CommitAnswers, error) {
	var answers CommitAnswers

	answers.Scope = defaultScope
	answers.ConfirmCommit = true
	answers.ConfirmPush = false

	fields := []huh.Field{
		// 1. Pemilihan Tipe Commit
		huh.NewSelect[string]().
			Title("Tipe Commit").
			Description("Pilih kategori perubahan kode Anda:").
			Options(
				huh.NewOption("feat (Fitur Baru)", "feat"),
				huh.NewOption("fix (Perbaikan Bug)", "fix"),
				huh.NewOption("docs (Dokumentasi)", "docs"),
				huh.NewOption("style (Format/Gaya Kode)", "style"),
				huh.NewOption("refactor (Restrukturisasi)", "refactor"),
				huh.NewOption("perf (Performa)", "perf"),
				huh.NewOption("test (Unit Test)", "test"),
				huh.NewOption("build (Build System/Deps)", "build"),
				huh.NewOption("ci (CI/CD Config)", "ci"),
				huh.NewOption("chore (Tugas Lain)", "chore"),
			).
			Value(&answers.Type),

		// 2. Scope Input
		huh.NewInput().
			Title("Scope Commit").
			Description("Modul/bagian kode yang diubah (misal: auth, db):").
			Value(&answers.Scope).
			Placeholder("opsional"),

		// 3. Subject Input
		huh.NewInput().
			Title("Deskripsi Singkat (Subject)").
			Description("Ringkasan singkat apa yang diubah (wajib):").
			Value(&answers.Subject).
			Placeholder("Tulis deskripsi commit...").
			Validate(func(str string) error {
				trimmed := strings.TrimSpace(str)

				if len(trimmed) == 0 {
					return fmt.Errorf("deskripsi wajib diisi!")
				}

				if len(trimmed) > 72 {
					return fmt.Errorf(
						"deskripsi terlalu panjang (maksimal 72 karakter)!",
					)
				}

				return nil
			}),

		// 4. Body Input
		huh.NewInput().
			Title("Deskripsi Detail (Body)").
			Description("Penjelasan lebih mendalam tentang perubahan (opsional):").
			Value(&answers.Body).
			Placeholder("opsional"),

		// 5. Konfirmasi Commit
		huh.NewConfirm().
			Title("Apakah Anda ingin membuat commit dengan pesan ini?").
			Value(&answers.ConfirmCommit),
	}

	// 6. Konfirmasi Push
	if branchName != "" {
		fields = append(
			fields,
			huh.NewConfirm().
				Title(
					fmt.Sprintf(
						"Apakah Anda ingin langsung melakukan push ke remote repository dengan branch '%s'?",
						branchName,
					),
				).
				Value(&answers.ConfirmPush),
		)
	}

	form := huh.NewForm(
		huh.NewGroup(fields...),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return &answers, nil
}

// AskAICommitConfirmation meminta konfirmasi user terhadap commit message
// yang sebelumnya sudah ditampilkan oleh caller.
func AskAICommitConfirmation() (bool, error) {
	var confirm = true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Gunakan commit message ini?").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirm, nil
}

// AskPushConfirmation meminta konfirmasi user untuk melakukan push.
func AskPushConfirmation(branchName string) (bool, error) {
	if branchName == "" {
		return false, nil
	}

	var confirmPush bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(
					fmt.Sprintf(
						"Apakah Anda ingin langsung melakukan push ke remote repository dengan branch '%s'?",
						branchName,
					),
				).
				Value(&confirmPush),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirmPush, nil
}

// FormatCommitMessage merangkai jawaban user menjadi
// satu string dengan format Conventional Commits.
func (a *CommitAnswers) FormatCommitMessage() string {
	var msg strings.Builder

	// Header: type(scope): subject
	msg.WriteString(a.Type)

	if a.Scope != "" {
		msg.WriteString("(")
		msg.WriteString(a.Scope)
		msg.WriteString(")")
	}

	msg.WriteString(": ")
	msg.WriteString(a.Subject)

	// Body jika diisi
	if strings.TrimSpace(a.Body) != "" {
		msg.WriteString("\n\n")
		msg.WriteString(strings.TrimSpace(a.Body))
	}

	return msg.String()
}
