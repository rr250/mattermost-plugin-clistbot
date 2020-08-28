package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/robfig/cron"
)

type Plugin struct {
	plugin.MattermostPlugin
	botUserID         string
	configurationLock sync.RWMutex
	configuration     *configuration
	cron              *cron.Cron
}
