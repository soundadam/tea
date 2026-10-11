package tui

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	appkit "github.com/soundadam/appkit/go"
	updatetui "github.com/soundadam/appkit/go/tui"

	"github.com/soundadam/tea/internal/app"
)

// updates holds the update prompt for the whole client session. The menu is a
// new bubbletea program after every action, so the prompt's state and its
// background work live here rather than in one program: a check or download
// that finishes while an action runs is applied all the same, and the next
// menu shows the result.
type updates struct {
	mu      sync.Mutex
	prompt  prompt
	program *tea.Program
}

// prompt is the part of appkit's updater the client uses.
type prompt interface {
	Update(tea.Msg) (prompt, tea.Cmd)
	View() string
}

type appkitPrompt struct{ updatetui.Updater }

func (p appkitPrompt) Update(msg tea.Msg) (prompt, tea.Cmd) {
	u, cmd := p.Updater.Update(msg)
	return appkitPrompt{u}, cmd
}

// updatesRefreshMsg asks the running menu to redraw after the prompt changed.
type updatesRefreshMsg struct{}

func newUpdates(info appkit.App) *updates {
	updater := updatetui.New(info)
	u := &updates{prompt: appkitPrompt{updater}}
	u.run(updater.Init())
	return u
}

// run executes one of the updater's commands off the UI and feeds its result
// back to the updater, whichever menu is on screen by then.
func (u *updates) run(cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	go func() {
		msg := cmd()
		u.mu.Lock()
		var next tea.Cmd
		u.prompt, next = u.prompt.Update(msg)
		program := u.program
		u.mu.Unlock()
		u.run(next)
		if program != nil {
			program.Send(updatesRefreshMsg{})
		}
	}()
}

func (u *updates) key(msg tea.KeyMsg) {
	u.mu.Lock()
	var cmd tea.Cmd
	u.prompt, cmd = u.prompt.Update(msg)
	u.mu.Unlock()
	u.run(cmd)
}

func (u *updates) view() string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.prompt.View()
}

func (u *updates) attach(program *tea.Program) {
	u.mu.Lock()
	u.program = program
	u.mu.Unlock()
}

// menuModel is the menu form with the update prompt under it.
type menuModel struct {
	form    *huh.Form
	updates *updates
}

func (m menuModel) Init() tea.Cmd { return m.form.Init() }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		m.updates.key(key)
	}
	model, cmd := m.form.Update(msg)
	if form, ok := model.(*huh.Form); ok {
		m.form = form
	}
	return m, cmd
}

func (m menuModel) View() string {
	view := m.form.View()
	if line := m.updates.view(); line != "" {
		view += "\n" + line
	}
	return view
}

// runMenu runs the menu form like runForm does, with the update prompt.
func runMenu(application app.App, form *huh.Form, u *updates) error {
	form = form.
		WithTheme(clientTheme()).
		WithKeyMap(clientKeyMap()).
		WithShowHelp(false)
	form.SubmitCmd = tea.Quit
	form.CancelCmd = tea.Interrupt
	program := tea.NewProgram(menuModel{form: form, updates: u}, programOptions(application)...)
	u.attach(program)
	_, err := program.Run()
	u.attach(nil)
	if form.State == huh.StateAborted || err == tea.ErrInterrupted {
		return huh.ErrUserAborted
	}
	return err
}
