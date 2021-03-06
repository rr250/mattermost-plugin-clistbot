package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) SendDailyContests() {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now()
	now = now.In(loc)
	nowStr := now.Format("2006-01-02")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://clist.by/api/v1/contest/?limit=20&offset=0&order_by=-start&start__gt="+nowStr+"T00:00:00&start__lt="+nowStr+"T23:59:59", nil)
	req.Header.Set("Authorization", "ApiKey rrrishabh7:5c0350a990cc6bb5dff68626e15e5f7f33a346f6")
	res, err := client.Do(req)

	if err != nil {
		p.API.LogError(err.Error())
		return
	}
	var body map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		p.API.LogError(err.Error())
		return
	}
	res.Body.Close()
	attachments := []*model.SlackAttachment{}
	objects := body["objects"].([]interface{})
	for _, object := range objects {
		object1 := object.(map[string]interface{})
		attachment := &model.SlackAttachment{}
		attachment.Title = object1["event"].(string)
		attachment.TitleLink = object1["href"].(string)
		attachment.AuthorName = object1["resource"].(map[string]interface{})["name"].(string)
		layout := "2006-01-02T15:04:05"
		startValue := object1["start"].(string)
		start, err1 := time.Parse(layout, startValue)
		if err1 != nil {
			p.API.LogError(err1.Error())
			return
		}
		start = start.In(loc)
		endValue := object1["end"].(string)
		end, err1 := time.Parse(layout, endValue)
		if err1 != nil {
			p.API.LogError(err1.Error())
			return
		}
		end = end.In(loc)
		attachment.Text = "Start : " + start.Format(time.RFC822) + " \n" + "End : " + end.Format(time.RFC822)
		attachments = append(attachments, attachment)
	}
	configuration := p.getConfiguration()
	for _, channelID := range strings.Split(strings.Trim(configuration.ChannelIDList, " "), ",") {
		postModel := &model.Post{
			UserId:    p.botUserID,
			ChannelId: channelID,
			Message:   "Challenges starting today",
			Props: model.StringInterface{
				"attachments": attachments,
			},
		}
		if len(attachments) == 0 {
			postModel.Message = "No Challenges today"
		}
		_, err2 := p.API.CreatePost(postModel)
		if err2 != nil {
			p.API.LogError(err2.Error())
			return
		}
	}

}
