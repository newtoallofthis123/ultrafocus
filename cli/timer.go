package cli

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
)

func GetTimerModel(duration time.Duration) timer.Model {
	t := timer.New(duration)

	return t
}
