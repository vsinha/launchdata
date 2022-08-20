package bubble

import (
	"fmt"
	"os"

	"launchdata/cli"
	"launchdata/config"
	"launchdata/parse"
	"launchdata/slices"

	"launchdata/list"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
)

// const listHeight = 14

var (
	appStyle = lipgloss.NewStyle().Padding(1, 0)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	viewportStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(0, 1)

	listWrapperStyle = lipgloss.NewStyle().Width(80).BorderStyle(lipgloss.HiddenBorder())
)

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	switchFocus      key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		switchFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch focus"),
		),
	}
}

type model struct {
	items        []MyItem
	list         list.Model[MyItem]
	listFocused  bool
	viewport     viewport.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func newModel(year int) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)
	entries, err := parse.LoadLaunchDataFromFile(fmt.Sprintf("./data/launchdata-%d.json", year))
	if err != nil {
		panic(err)
	}

	var items []MyItem
	for _, r := range entries.OrbitalFlights {
		items = append(items, MyItem{data: r})
	}
	slices.Reverse(items)

	width, height, err := terminal.GetSize(0)
	height = height - 5
	listWidth := int(float32(width) * 0.5)

	delegate := newItemDelegate(delegateKeys)
	l := list.New[MyItem](items, delegate, listWidth, height)
	l.Title = fmt.Sprintf("Launches in %d", year)
	l.Styles.Title = titleStyle
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	viewport := viewport.New(width-listWidth, height)
	viewport.Style = viewportStyle

	return model{
		items:        items,
		list:         l,
		listFocused:  true,
		viewport:     viewport,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			cmds = append(cmds, cmd)

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.switchFocus):
			m.listFocused = !m.listFocused
		}
	}

	if m.listFocused {
		list, cmd := m.list.Update(msg)
		m.list = list
		cmds = append(cmds, cmd)
	}

	if i := m.list.SelectedItem(); i != nil {
		m.viewport.SetContent(i.Render(80))
	} else {
		m.viewport.SetContent("")
	}

	if !m.listFocused {
		viewport, cmd := m.viewport.Update(msg)
		m.viewport = viewport
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.listFocused {
		listWrapperStyle.BorderStyle(lipgloss.RoundedBorder())
		m.viewport.Style.BorderStyle(lipgloss.HiddenBorder())
	} else {
		listWrapperStyle.BorderStyle(lipgloss.HiddenBorder())
		m.viewport.Style.BorderStyle(lipgloss.RoundedBorder())
	}

	list := listWrapperStyle.Render(m.list.View())
	view := m.viewport.View()

	return lipgloss.JoinHorizontal(lipgloss.Left, list, view)
}

func Run(config *config.Config, year int) {
	cli.ClearScreen()

	p := tea.NewProgram(newModel(year), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
