package bubble

import (
	"fmt"

	"launchdata/list"
	"launchdata/parse"

	"github.com/charmbracelet/lipgloss"
)

type MyItem struct {
	data parse.RocketData
}

func (i MyItem) Title() string {
	return fmt.Sprintf("%s (%s)",
		i.data.Rocket,
		i.data.FlightNumber)
}

func (i MyItem) Description() string {
	return fmt.Sprintf("%v, %s, %s",
		i.data.Timestamp.DateString(),
		i.data.LaunchServiceProvider,
		i.data.LaunchSite)
}

func (i MyItem) FilterValue() string {
	return fmt.Sprintf("%v", i.data)
}

func (i MyItem) Render(width int) string {
	return lipgloss.NewStyle().Width(80).Render(i.data.Render())
}

var _ list.Item = (*MyItem)(nil)
