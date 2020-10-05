package main

import (
	"sync"

	"github.com/mattermost/mattermost-plugin-api/cluster"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	botUserID         string
	configurationLock sync.RWMutex
	configuration     *configuration
	backgroundJob     *cluster.Job
}
