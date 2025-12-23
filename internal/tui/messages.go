package tui

import (
	"fmt"
	"github.com/karol-broda/snitch/internal/collector"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type dataMsg struct {
	connections []collector.Connection
}

type errMsg struct {
	err error
}

type killResultMsg struct {
	pid     int
	process string
	success bool
	err     error
}

type clearStatusMsg struct{}

func (m model) tick() tea.Cmd {
	return tea.Tick(m.interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) fetchData() tea.Cmd {
	return func() tea.Msg {
		conns, err := collector.GetConnections()
		if err != nil {
			return errMsg{err}
		}
		return dataMsg{connections: conns}
	}
}

func killProcess(pid int, process string) tea.Cmd {
	return func() tea.Msg {
		if pid <= 0 {
			return killResultMsg{
				pid:     pid,
				process: process,
				success: false,
				err:     fmt.Errorf("invalid pid"),
			}
		}

		// send SIGTERM first (graceful shutdown)
		err := syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			return killResultMsg{
				pid:     pid,
				process: process,
				success: false,
				err:     err,
			}
		}

		return killResultMsg{
			pid:     pid,
			process: process,
			success: true,
			err:     nil,
		}
	}
}

func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

