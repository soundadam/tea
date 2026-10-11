package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"strings"

	"github.com/soundadam/tea/internal/power"
	"github.com/soundadam/tea/internal/shutdown"
)

func testMenu(u *updates) menuModel {
	var action string
	form := newMenuForm(power.Status{Observation: power.Off}, shutdown.Status{}, "", nil, &action)
	menu := menuModel{form: form, updates: u}
	menu.Init()
	return menu
}

// fakePrompt starts its work on "u" and shows whatever its last result said.
type fakePrompt struct {
	line    string
	release chan struct{}
}

type fakeDone string

func (p fakePrompt) Update(msg tea.Msg) (prompt, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "u" {
			p.line = "working"
			release := p.release
			return p, func() tea.Msg { <-release; return fakeDone("done") }
		}
	case fakeDone:
		p.line = string(msg)
	}
	return p, nil
}

func (p fakePrompt) View() string { return p.line }

func TestUpdateWorkFinishesWithoutTheMenuThatStartedIt(t *testing.T) {
	release := make(chan struct{})
	u := &updates{prompt: fakePrompt{release: release}}
	menu := testMenu(u)

	model, _ := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")})
	if got := u.view(); got != "working" {
		t.Fatalf("after the key the prompt shows %q", got)
	}
	if view := model.View(); !strings.HasSuffix(view, "\nworking") {
		t.Fatalf("the menu does not end with the prompt line:\n%s", view)
	}

	// The menu program has gone (an action is running); the work completes anyway.
	close(release)
	deadline := time.Now().Add(2 * time.Second)
	for u.view() != "done" {
		if time.Now().After(deadline) {
			t.Fatalf("the result never reached the prompt: %q", u.view())
		}
		time.Sleep(time.Millisecond)
	}
}

func TestMenuWithoutPromptLineIsTheForm(t *testing.T) {
	u := &updates{prompt: fakePrompt{}}
	menu := testMenu(u)
	if menu.View() != menu.form.View() {
		t.Fatal("an empty prompt adds lines to the menu")
	}
}
