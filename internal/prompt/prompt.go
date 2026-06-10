package prompt

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// CommitAnswers menyimpan jawaban dari user untuk setiap bagian Conventional Commits.
type CommitAnswers struct {
	Type      string
	Scope     string
	Subject   string
	Body      string
	HasFooter bool
	Footer    string
}

// AskQuestions memicu rangkaian form interaktif di terminal.
func AskQuestions() (*CommitAnswers, error) {
	var answers CommitAnswers

	// 1. Tipe Commit (pilihan ganda)
	typePrompt := &survey.Select{
		Message: "Pilih tipe commit:",
		Options: []string{
			"feat:     Fitur baru",
			"fix:      Perbaikan bug",
			"docs:     Perubahan dokumentasi",
			"style:    Perbaikan format kode (semicolon, spasi, dll)",
			"refactor: Restrukturisasi kode (bukan bugfix atau fitur)",
			"perf:     Peningkatan performa",
			"test:     Menambah atau memperbaiki unit test",
			"build:    Perubahan build system atau dependency (go.mod, webpack, dll)",
			"ci:       Perubahan konfigurasi CI/CD (GitHub Actions, dll)",
			"chore:    Tugas rutin lain (tanpa menyentuh source code utama)",
			"revert:   Membatalkan commit sebelumnya",
		},
		Default: "feat:     Fitur baru",
	}

	var selectedType string
	if err := survey.AskOne(typePrompt, &selectedType); err != nil {
		return nil, err
	}

	// Ambil bagian tipenya saja sebelum tanda titik dua (misal: "feat")
	answers.Type = strings.TrimSpace(strings.Split(selectedType, ":")[0])

	// 2. Scope (opsional, untuk memperjelas modul mana yang diubah)
	scopePrompt := &survey.Input{
		Message: "Masukkan scope commit (opsional, contoh: auth, database):",
	}
	if err := survey.AskOne(scopePrompt, &answers.Scope); err != nil {
		return nil, err
	}
	answers.Scope = strings.TrimSpace(answers.Scope)

	// 3. Subject (wajib, deskripsi singkat)
	subjectPrompt := &survey.Input{
		Message: "Tuliskan deskripsi singkat commit (wajib):",
	}
	if err := survey.AskOne(subjectPrompt, &answers.Subject, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}
	answers.Subject = strings.TrimSpace(answers.Subject)

	// 4. Body (opsional, penjelasan mendalam)
	bodyPrompt := &survey.Input{
		Message: "Tuliskan deskripsi detail/body commit (opsional):",
	}
	if err := survey.AskOne(bodyPrompt, &answers.Body); err != nil {
		return nil, err
	}
	answers.Body = strings.TrimSpace(answers.Body)

	// 5. Konfirmasi apakah butuh footer (Breaking changes atau referensi isu/issue number)
	confirmPrompt := &survey.Confirm{
		Message: "Apakah ada breaking changes atau isu (GitHub) yang ingin ditutup?",
		Default: false,
	}
	if err := survey.AskOne(confirmPrompt, &answers.HasFooter); err != nil {
		return nil, err
	}

	// 6. Jika butuh footer, isi deskripsinya
	if answers.HasFooter {
		footerPrompt := &survey.Input{
			Message: "Tuliskan isi footer (contoh: Closes #12, BREAKING CHANGE: ganti nama parameter x):",
		}
		if err := survey.AskOne(footerPrompt, &answers.Footer); err != nil {
			return nil, err
		}
		answers.Footer = strings.TrimSpace(answers.Footer)
	}

	return &answers, nil
}

// FormatCommitMessage merangkai jawaban user menjadi satu string format Conventional Commits.
func (a *CommitAnswers) FormatCommitMessage() string {
	var msg strings.Builder

	// Menulis header commit: type(scope): subject
	msg.WriteString(a.Type)
	if a.Scope != "" {
		msg.WriteString("(" + a.Scope + ")")
	}
	msg.WriteString(": ")
	msg.WriteString(a.Subject)

	// Menulis body jika ada
	if a.Body != "" {
		msg.WriteString("\n\n")
		msg.WriteString(a.Body)
	}

	// Menulis footer jika ada
	if a.HasFooter && a.Footer != "" {
		msg.WriteString("\n\n")
		msg.WriteString(a.Footer)
	}

	return msg.String()
}
