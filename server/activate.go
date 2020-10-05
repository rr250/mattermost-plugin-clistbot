package main

import (
	"time"

	"github.com/mattermost/mattermost-plugin-api/cluster"
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
	job, cronErr := cluster.Schedule(
		p.API,
		"BackgroundJob",
		cluster.MakeWaitForRoundedInterval(24*time.Hour),
		p.SendDailyContests,
	)
	if cronErr != nil {
		return errors.Wrap(cronErr, "failed to schedule background job")
	}
	p.backgroundJob = job
	return nil
}

func (p *Plugin) ensureBotExists() (string, error) {
	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}

	return p.Helpers.EnsureBot(bot)
}
