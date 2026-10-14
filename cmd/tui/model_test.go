package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSplitUrls(t *testing.T) {
	in := "https://youtu.be/a\nhttps://www.youtube.com/watch?v=b, https://music.youtube.com/watch?v=c   https://youtu.be/a"
	got := splitUrls(in)
	want := []string{
		"https://youtu.be/a",
		"https://www.youtube.com/watch?v=b",
		"https://music.youtube.com/watch?v=c",
	}
	if len(got) != len(want) {
		t.Fatalf("splitUrls length = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("splitUrls[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestJobFileName(t *testing.T) {
	cases := map[string]string{
		"downloads/My Song.mp3": "My Song.mp3",
		`downloads\Track.mp3`:   "Track.mp3",
		"":                      "",
		"plain.mp3":             "plain.mp3",
	}
	for in, want := range cases {
		if got := (Job{Output: in}).FileName(); got != want {
			t.Errorf("Job{Output:%q}.FileName() = %q, want %q", in, got, want)
		}
	}
}

func TestItoa(t *testing.T) {
	if itoa(0) != "0" || itoa(42) != "42" || itoa(-7) != "-7" {
		t.Errorf("itoa produced unexpected output: %q %q %q", itoa(0), itoa(42), itoa(-7))
	}
}

func TestModelStartsOnLogin(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	if m.screen != screenLogin {
		t.Fatalf("expected initial screen to be login")
	}
	view := m.View()
	if !strings.Contains(view, "traer") {
		t.Errorf("login view missing app title:\n%s", view)
	}
}

func TestLoginToMainTransition(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.identifier.SetValue("alice")
	m.password.SetValue("secret")

	next, _ := m.submitLogin()
	// submitLogin returns a command; simulate its successful result.
	updated, _ := next.Update(loginResultMsg{err: nil})
	mm := updated.(model)
	if mm.screen != screenMain {
		t.Fatalf("expected screen to switch to main after login, got %v", mm.screen)
	}
	if !strings.Contains(mm.View(), "Music") {
		t.Errorf("main view missing music heading:\n%s", mm.View())
	}
}

func TestSubmitMusicValidatesInput(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.screen = screenMain
	m.urls.SetValue("   ") // only whitespace

	next, _ := m.submitMusic()
	mm := next.(model)
	if mm.err == "" {
		t.Errorf("expected an error when no URLs are provided")
	}
}

func TestMergeJobsDedupAndPreserve(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.screen = screenMain

	// Seed with a completed job that will no longer appear in "active".
	m.mergeJobs([]Job{{ID: "done-1", Status: "done", Uploaded: true}})

	// Active refresh returns a running job plus an update to nothing new.
	m.mergeJobs([]Job{{ID: "run-1", Status: "running", Progress: 40}})

	if len(m.jobs) != 2 {
		t.Fatalf("expected 2 jobs after merge, got %d", len(m.jobs))
	}

	// Refresh again with an update to run-1; done-1 must be preserved.
	m.mergeJobs([]Job{{ID: "run-1", Status: "done", Uploaded: true}})
	if len(m.jobs) != 2 {
		t.Fatalf("expected 2 jobs after second merge, got %d", len(m.jobs))
	}

	var sawDone, sawRun bool
	for _, j := range m.jobs {
		switch j.ID {
		case "done-1":
			sawDone = true
		case "run-1":
			sawRun = true
			if !j.Uploaded {
				t.Errorf("run-1 should have been updated to uploaded")
			}
		}
	}
	if !sawDone || !sawRun {
		t.Errorf("merge lost a job: sawDone=%v sawRun=%v", sawDone, sawRun)
	}
}

func TestHandleMainKeyToggleFocus(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.screen = screenMain
	if m.mainFocus != focusURLs {
		t.Fatalf("expected initial main focus to be URLs")
	}
	next, _ := m.handleMainKey(tea.KeyMsg{Type: tea.KeyTab})
	if next.(model).mainFocus != focusJobs {
		t.Errorf("expected tab to move focus to jobs")
	}
	next2, _ := next.(model).handleMainKey(tea.KeyMsg{Type: tea.KeyTab})
	if next2.(model).mainFocus != focusURLs {
		t.Errorf("expected second tab to move focus back to URLs")
	}
}

func TestProgressBarBounds(t *testing.T) {
	if got := progressBar(-10, 10); strings.Contains(got, "█") {
		t.Errorf("progressBar(-10) should be empty, got %q", got)
	}
	if got := progressBar(150, 10); len([]rune(got)) != 10 {
		t.Errorf("progressBar(150) width = %d, want 10", len([]rune(got)))
	}
	if got := progressBar(50, 10); strings.Count(got, "█") != 5 {
		t.Errorf("progressBar(50) filled = %d, want 5", strings.Count(got, "█"))
	}
}
