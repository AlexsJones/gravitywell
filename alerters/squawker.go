package alerters

import (
	"fmt"
	"log"
	"sync"

	notifier "github.com/AlexsJones/squawker"
	"github.com/AlexsJones/squawker/services/slack"
)

var (
	instance     *notifier.Manager
	once         sync.Once
	slackToken   = "xoxp-2310897947-180906251303-376634262709-170040c3c539cc7f804ea840409293d3"
	slackChannel = "team_disco"
)

//GetNotifier ....
func GetNotifier() *notifier.Manager {
	once.Do(func() {

		notifierManager := notifier.NewManager(func(notifier notifier.INotifier, err error) {

			log.Printf(fmt.Sprintf("Notifier %s error: %s\n", notifier.GetName(), err.Error()))
		})

		notifierManager.AddNotifier(&slack.Notifier{ClientToken: slackToken, Channels: []string{slackChannel}})

		instance = notifierManager
	})
	return instance
}
