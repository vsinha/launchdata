package bubble

import (
	"fmt"

	"launchdata/list"
	"launchdata/parse"
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
		i.data.Timestamp.Format("01-02-2006"),
		i.data.LaunchServiceProvider,
		i.data.LaunchSite)
}

func (i MyItem) FilterValue() string {
	return fmt.Sprintf("%v", i.data)
}

var _ list.Item = (*MyItem)(nil)
