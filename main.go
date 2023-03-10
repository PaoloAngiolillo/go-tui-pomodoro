package main

// A simple example that shows how to render a progress bar in a "pure"
// fashion. In this example we bump the progress by 25% every second,
// maintaining the progress state on our top level model using the progress bar
// model's ViewAs method only for rendering.
//
// The signature for ViewAs is:
//
//     func (m Model) ViewAs(percent float64) string
//
// So it takes a float between 0 and 1, and renders the progress bar
// accordingly. When using the progress bar in this "pure" fashion and there's
// no need to call an Update method.
//
// The progress bar is also able to animate itself, however. For details see
// the progress-animated example.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding       = 3
	paddingMiddle = 40
	maxWidth      = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	// prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	percent        float64
	progress       progress.Model
	secondsPassed  int
	startTime      time.Time
	endTime        time.Time
	timerRemaining time.Duration
	timerDuration  time.Duration
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+s", "s":
			m.timerDuration = 1500
			m.startTime = time.Now()
			m.endTime = m.startTime.Add(m.timerDuration * time.Second)
			return m, tickCmd()
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
			m.timerRemaining = 0
		}
		return m, nil

	case tickMsg:
		m.timerRemaining = time.Until(m.endTime).Round(1 * time.Second)
		m.secondsPassed += 1
		m.percent = float64(m.secondsPassed) / float64(m.timerDuration)
		if m.percent > 1.0 {
			m.percent = 1.0
			return m, tea.Println("Pomodoro Timer done!")
		}
		return m, tickCmd()

	default:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	padMiddle := strings.Repeat(" ", paddingMiddle)
	return "\n" +
		padMiddle + "\n\n" +
		padMiddle + m.timerRemaining.String() + "\n\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press ctrl+c or q to quit, Press ctrl+s or s to start. ") + "\n\n"
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	return model{
		progress:       prog,
		timerDuration:  120,
		timerRemaining: 120,
		percent:        0,
	}
}
