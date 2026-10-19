package main

import (
	"strconv"
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

func TestIsYouTubeURL(t *testing.T) {
	valid := []string{
		"https://youtu.be/dQw4w9WgXcQ",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://music.youtube.com/watch?v=dQw4w9WgXcQ",
		"youtube.com/watch?v=abc",
	}
	for _, u := range valid {
		if !IsYouTubeURL(u) {
			t.Errorf("IsYouTubeURL(%q) = false, want true", u)
		}
	}
	invalid := []string{
		"https://example.com/video.mp4",
		"not a url",
		"https://vimeo.com/12345",
	}
	for _, u := range invalid {
		if IsYouTubeURL(u) {
			t.Errorf("IsYouTubeURL(%q) = true, want false", u)
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
	m.email.SetValue("alice@example.com")
	m.password.SetValue("secret")

	next, _ := m.submitLogin()
	// submitLogin returns a command; simulate its successful result.
	updated, _ := next.Update(loginResultMsg{userID: "u1", err: nil})
	mm := updated.(model)
	if mm.screen != screenMain {
		t.Fatalf("expected screen to switch to main after login, got %v", mm.screen)
	}
	if !strings.Contains(mm.View(), "Jobs") {
		t.Errorf("main view missing jobs heading:\n%s", mm.View())
	}
}

func TestSubmitJobsValidatesInput(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.screen = screenMain

	// Only whitespace.
	m.urls.SetValue("   ")
	next, _ := m.submitJobs()
	if mm := next.(model); mm.err == "" {
		t.Errorf("expected an error when no URLs are provided")
	}

	// A non-YouTube URL must be rejected before hitting the backend.
	m.urls.SetValue("https://example.com/video.mp4")
	next, _ = m.submitJobs()
	if mm := next.(model); mm.err == "" {
		t.Errorf("expected an error for a non-YouTube URL")
	}
}

func TestSetJobsCapsAndClamps(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.screen = screenMain

	// More jobs than the cap: only the newest maxTrackedJobs are kept.
	var many []Job
	for i := 0; i < maxTrackedJobs+10; i++ {
		many = append(many, Job{ID: "job-" + strconv.Itoa(i), Status: statusQueued})
	}
	m.setJobs(many)
	if len(m.jobs) != maxTrackedJobs {
		t.Fatalf("expected %d jobs after cap, got %d", maxTrackedJobs, len(m.jobs))
	}
	if m.jobs[0].ID != "job-10" {
		t.Errorf("expected newest jobs to be kept, first = %q", m.jobs[0].ID)
	}

	// A stale selection index must be clamped when the list shrinks.
	m.jobIndex = 40
	m.setJobs([]Job{{ID: "only", Status: statusCompleted}})
	if m.jobIndex != 0 {
		t.Errorf("expected jobIndex clamped to 0, got %d", m.jobIndex)
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

func TestProcessTypeCycles(t *testing.T) {
	m := newModel(NewClient("http://example.test"))
	m.screen = screenMain

	if m.currentProcessType() != "audio" {
		t.Fatalf("expected default process type audio, got %q", m.currentProcessType())
	}

	// ctrl+p advances to the next backend-supported process type.
	next, _ := m.handleMainKey(tea.KeyMsg{Type: tea.KeyCtrlP})
	mm := next.(model)
	if mm.currentProcessType() != "video" {
		t.Fatalf("expected video after one cycle, got %q", mm.currentProcessType())
	}

	// Cycling wraps around to the start after the last type.
	for i := 1; i < len(processTypes); i++ {
		next, _ = mm.handleMainKey(tea.KeyMsg{Type: tea.KeyCtrlP})
		mm = next.(model)
	}
	if mm.currentProcessType() != processTypes[0] {
		t.Fatalf("expected wrap-around to %q, got %q", processTypes[0], mm.currentProcessType())
	}
}
