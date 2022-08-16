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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
)

type item struct {
	data parse.RocketData
}

func (i item) Title() string {
	return fmt.Sprintf("%s (%s)",
		i.data.Rocket,
		i.data.FlightNumber)
}

func (i item) Description() string {
	return fmt.Sprintf("%v, %s, %s",
		i.data.Timestamp.Format("01-02-2006"),
		i.data.LaunchServiceProvider,
		i.data.LaunchSite)
}

func (i item) FilterValue() string {
	return fmt.Sprintf("%v", i.data)
}

var _ list.Item = (*item)(nil)

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	nextPage         key.Binding
	prevPage         key.Binding
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
		nextPage: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "next page"),
		),
		prevPage: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("Ctrl+U", "previous page"),
		),
	}
}

type model struct {
	list         list.Model
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

	var items []list.Item
	for _, r := range entries.OrbitalFlights {
		items = append(items, item{data: r})
	}
	slices.Reverse(items)

	delegate := newItemDelegate(delegateKeys)
	l := list.New(items, delegate, 0, 0)
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
			listKeys.nextPage,
			listKeys.prevPage,
		}
	}

	return model{
		list:         l,
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

		case key.Matches(msg, m.keys.prevPage):
			m.list.Paginator.PrevPage()

		case key.Matches(msg, m.keys.nextPage):
			m.list.Paginator.NextPage()
		}
	}

	newlist, cmd := m.list.Update(msg)
	m.list = newlist
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
	// if m.choice != "" {
	// 	return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	// }
	// if m.quitting {
	// 	return quitTextStyle.Render("Not hungry? That’s cool.")
	// }
	// return "\n" + m.list.View()
}

func Run(config *config.Config, year int) {
	cli.ClearScreen()

	p := tea.NewProgram(newModel(year), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
