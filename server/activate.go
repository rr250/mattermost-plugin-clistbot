package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	botUserName    = "contest_bot"
	botDisplayName = "Contest Finder Bot"
)

// OnActivate register the plugin command
func (p *Plugin) OnActivate() error {
	botUserID, err := p.ensureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot user")
	}
	p.botUserID = botUserID
	p.cron = p.InitCRON()
	p.cron.Start()
	return nil
}

func (p *Plugin) ensureBotExists() (string, error) {
	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}

	return p.Helpers.EnsureBot(bot)
}
