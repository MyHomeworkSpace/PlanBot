package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v31/github"
	"github.com/slack-go/slack"
)

var slackAPI *slack.Client
var config = Config{}

func main() {
	configure(&config)

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, config.GitHub.AppID, config.GitHub.InstallationID, config.GitHub.PrivateKey)
	if err != nil {
		panic(err)
	}

	timeUntilUnassigned := time.Hour * 24 * time.Duration(config.Unassign.DaysUntilUnassign)
	timeUntilNotified := time.Hour * 24 * time.Duration(config.Unassign.DaysUntilWarning)

	gh := github.NewClient(&http.Client{Transport: itr})
	slackAPI = slack.New(config.Slack.Token)

	if err != nil {
		panic(err)
	}

	opts := &github.IssueListByRepoOptions{Assignee: "*", State: "open"}

	issues, _, err := gh.Issues.ListByRepo(context.Background(), config.GitHub.Organization, config.GitHub.Repo, opts)
	if err != nil {
		logError("[FATAL] Getting issues", err)
		panic(err)
	}

	for _, issue := range issues {
		assignees := []string{}
		for _, assignee := range issue.Assignees {
			assignees = append(assignees, *assignee.Login)
		}
		if issue.UpdatedAt.Add(timeUntilUnassigned).Before(time.Now()) {
			// the issue hasn't been updated in a week
			message := "Hi! I've removed @" + strings.Join(assignees, ", @") + " from this issue due to a lack of activity. If this is a mistake, feel free to reassign yourself. If you would like to work on this issue, feel free to assign yourself."
			_, _, err = gh.Issues.RemoveAssignees(context.Background(), config.GitHub.Organization, config.GitHub.Repo, *issue.Number, assignees)
			if err != nil {
				logError("Removing issue assignees", err)
				continue
			}
			_, _, err = gh.Issues.CreateComment(context.Background(), config.GitHub.Organization, config.GitHub.Repo, *issue.Number, &github.IssueComment{
				Body: github.String(message),
			})
			if err != nil {
				logError("Creating issue comment", err)
			}
		} else if issue.UpdatedAt.Add(timeUntilNotified).Before(time.Now()) {
			for _, assignee := range assignees {
				message := "Hi! This is a warning that due to a lack of activity on <" + *issue.HTMLURL + "|Issue #" + strconv.Itoa(*issue.Number) + ">, you will be unassigned soon. Comment on the issue if you want to stay assigned (A simple \"bump\" will suffice)."
				_, _, _, err := slackAPI.SendMessage(config.Users[assignee], slack.MsgOptionText(message, false))
				if err != nil {
					logError("Posting message to slack", err)
				}
			}
		}
	}
}

func logError(desc string, err error) {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, false)
	stackTrace := string(buf[0:stackSize])

	log.Println("======================================")

	log.Printf("Error occurred while '%s'!", desc)
	errDesc := ""
	if err != nil {
		errDesc = err.Error()
	} else {
		errDesc = "(err == nil)"
	}
	log.Println(errDesc)
	log.Println(stackTrace)

	log.Println("======================================")

	title := fmt.Sprintf("An error occurred - %s", desc)
	_, _, _, postingErr := slackAPI.SendMessage(config.Slack.ErrlogChannel, slack.MsgOptionAttachments(slack.Attachment{
		Fallback:   title,
		Color:      "danger",
		Title:      title,
		Text:       "```" + stackTrace + "```",
		MarkdownIn: []string{"fields"},
		Fields: []slack.AttachmentField{
			{
				Title: "Host",
				Value: "PlanBot :robot_face:",
				Short: true,
			},
			{
				Title: "Message",
				Value: err.Error(),
				Short: true,
			},
		},
	}))
	if postingErr != nil {
		log.Println("An error occurred while reporting the previous error to Slack")
		log.Println(postingErr)
	}
}
