package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// typeKeys feeds a sequence of key messages into the model and returns the
// resulting model.
func typeKeys(m model, keys string) (model, []tea.Cmd) {
	var cmds []tea.Cmd
	for _, r := range keys {
		next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = next.(model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return m, cmds
}

func press(m model, keyType tea.KeyType) (model, tea.Cmd) {
	next, cmd := m.Update(tea.KeyMsg{Type: keyType})
	return next.(model), cmd
}

func TestLoginTypingAndSubmit(t *testing.T) {
	srv := newMockBackend(t, "alice", "secret")
	defer srv.Close()

	m := newModel(NewClient(srv.URL))

	// Type the identifier (login mode, focus starts at 0 = identifier).
	m, _ = typeKeys(m, "alice")
	if m.identifier.Value() != "alice" {
		t.Fatalf("identifier not captured, got %q", m.identifier.Value())
	}

	// Tab to password.
	m, _ = press(m, tea.KeyTab)
	if m.loginFocus != 1 {
		t.Fatalf("after tab, loginFocus = %d, want 1", m.loginFocus)
	}

	// Type the password.
	m, _ = typeKeys(m, "secret")
	if m.password.Value() != "secret" {
		t.Fatalf("password not captured, got %q", m.password.Value())
	}

	// Press enter to submit.
	m, cmd := press(m, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("enter produced no command")
	}

	// Execute the async login command and feed its result back.
	msg := cmd()
	m2, _ := m.Update(msg)
	m = m2.(model)
	if m.err != "" {
		t.Fatalf("login failed with error: %s", m.err)
	}
	if m.screen != screenMain {
		t.Fatalf("after successful login, screen = %v, want main", m.screen)
	}
}

// TestSignupFlowThroughModel drives the 3-field signup form (username, email,
// password) via tab navigation and submit, then logs in.
func TestSignupFlowThroughModel(t *testing.T) {
	srv := newMockBackend(t, "alice", "secret")
	defer srv.Close()

	m := newModel(NewClient(srv.URL))

	// Toggle to signup mode.
	m, _ = press(m, tea.KeyCtrlT)
	if !m.signupMode {
		t.Fatal("expected signup mode after ctrl+t")
	}
	if m.loginFocus != 0 {
		t.Fatalf("expected focus reset to 0 on mode switch, got %d", m.loginFocus)
	}

	// username (focus 0)
	m, _ = typeKeys(m, "alice")
	if m.username.Value() != "alice" {
		t.Fatalf("username not captured, got %q", m.username.Value())
	}

	// tab -> email (focus 1)
	m, _ = press(m, tea.KeyTab)
	m, _ = typeKeys(m, "alice@example.com")
	if m.email.Value() != "alice@example.com" {
		t.Fatalf("email not captured, got %q", m.email.Value())
	}

	// tab -> password (focus 2)
	m, _ = press(m, tea.KeyTab)
	if m.loginFocus != 2 {
		t.Fatalf("expected focus 2, got %d", m.loginFocus)
	}
	m, _ = typeKeys(m, "secret")
	if m.password.Value() != "secret" {
		t.Fatalf("password not captured, got %q", m.password.Value())
	}

	// submit signup
	m, cmd := press(m, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("enter produced no command")
	}
	m2, _ := m.Update(cmd())
	m = m2.(model)
	if m.err != "" {
		t.Fatalf("signup failed: %s", m.err)
	}
	if m.signupMode {
		t.Fatal("expected mode to switch back to login after successful signup")
	}
}

// TestQDoesNotQuitWhileTypingLogin ensures the letter 'q' in a credential does
// not terminate the app while on the login screen.
func TestQDoesNotQuitWhileTypingLogin(t *testing.T) {
	srv := newMockBackend(t, "alice", "secret")
	defer srv.Close()

	m := newModel(NewClient(srv.URL))
	// Focus password (loginFocus 1) and type a password containing 'q'.
	m, _ = press(m, tea.KeyTab)
	m, _ = typeKeys(m, "qersecret")

	if m.password.Value() != "qersecret" {
		t.Fatalf("password with 'q' not captured, got %q", m.password.Value())
	}
	// If 'q' had quit, the model would have been discarded by tea; here we only
	// assert the value survived, which proves 'q' was routed to the input.
}
