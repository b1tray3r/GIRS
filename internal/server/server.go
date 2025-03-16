package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	gitea "code.gitea.io/sdk/gitea"
	"github.com/b1tray3r/go-openapi3/pkg/api"
	redmine "github.com/nixys/nxs-go-redmine/v5"
)

const (
	RedmineTrackerID      = "RedmineTrackerID"      // Definiton for Task ("Aufgabe")
	RedmineClosedStatusID = "RedmineClosedStatusID" // Definition for "Closed" status
)

// Server is the struct used to implement the API
type Server struct {
	Redmine           *redmine.Context
	RedmineProjectKey string

	RedmineConfig map[string]int64

	Gitea *gitea.Client

	Echo *echo.Echo
}

// GetHealth is the implementation of the API endpoint /health
func (s Server) GetHealth(ctx context.Context, request api.GetHealthRequestObject) (api.GetHealthResponseObject, error) {
	message := time.Now().Format(time.RFC3339)
	response := api.GetHealth200JSONResponse{
		HealthCheckResponseJSONResponse: api.HealthCheckResponseJSONResponse{
			Message: &message,
		},
	}

	return response, nil
}

// PostGiteaWebhook is the implementation of the API endpoint /gitea/webhook
func (s Server) PostGiteaWebhook(ctx context.Context, request api.PostGiteaWebhookRequestObject) (api.PostGiteaWebhookResponseObject, error) {
	message := ""

	switch *request.Body.Action {
	case "closed":
		message = "Success! Issue closed in Redmine."
		err := s.closeIssue(request)
		if err != nil {
			message = err.Error()
		}
	case "opened":
		issueID, commentID, err := s.createIssue(request)
		message = fmt.Sprintf("Success! Issue created with ID: %d. Comment created with ID: %d", issueID, commentID)
		if err != nil {
			message = err.Error()
		}
	}

	response := api.PostGiteaWebhook200JSONResponse{
		WebhookAcceptedResponseJSONResponse: api.WebhookAcceptedResponseJSONResponse{
			Message: &message,
		},
	}

	return response, nil
}

// closeIssue closes the issue in Redmine if the issue contains a comment with the redmine issue url
// The first issue found will be closed.
func (s Server) closeIssue(request api.PostGiteaWebhookRequestObject) error {
	comments, _, err := s.Gitea.ListIssueComments(
		*request.Body.Repository.Owner.Login,
		*request.Body.Repository.Name,
		*request.Body.Number,
		gitea.ListIssueCommentOptions{},
	)
	if err != nil {
		return fmt.Errorf("error getting comments: %w", err)
	}

	for _, comment := range comments {
		if strings.Contains(comment.Body, "IssueInRedmine: https://projects.sdzecom.de/issues/") {
			parts := strings.Split(comment.Body, "/")
			issueID := parts[len(parts)-1]
			issueIDInt64, err := strconv.ParseInt(issueID, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing issue id: %w", err)
			}
			closeState := s.RedmineConfig[RedmineClosedStatusID]
			note := "Closed by webhook"
			code, err := s.Redmine.IssueUpdate(issueIDInt64, redmine.IssueUpdate{
				Issue: redmine.IssueUpdateObject{
					Notes:    &note,
					StatusID: &closeState,
				},
			})
			if err != nil {
				return fmt.Errorf("error closing issue with error: %w", err)
			}

			if code >= 400 {
				return fmt.Errorf("error closing issue with code: %d", code)
			}

			return nil // only close the first issue found
		}
	}

	return nil
}

// createIssue creates an issue in Redmine and comments the redmine issue url on the gitea issue
func (s Server) createIssue(request api.PostGiteaWebhookRequestObject) (int64, int64, error) {
	// Get the DevOps project
	p, _, err := s.Redmine.ProjectSingleGet(s.RedmineProjectKey, redmine.ProjectSingleGetRequest{})
	if err != nil {
		return -1, -1, fmt.Errorf("error getting project with error: %w", err)
	}

	// Create the issue
	redmineTrackerID := s.RedmineConfig[RedmineTrackerID]
	issue, code, err := s.Redmine.IssueCreate(redmine.IssueCreate{
		Issue: redmine.IssueCreateObject{
			ProjectID:   p.ID,
			TrackerID:   &redmineTrackerID,
			Subject:     *request.Body.Issue.Title,
			Description: request.Body.Issue.Body,
		},
	})
	if err != nil {
		return -1, -1, fmt.Errorf("error creating issue with error: %w", err)
	}
	if code != 201 {
		return -1, -1, fmt.Errorf("error creating issue with unexpected code: %d", code)
	}

	// Comment the redmine issue url on the gitea issue
	comment, _, err := s.Gitea.CreateIssueComment(
		*request.Body.Repository.Owner.Login,
		*request.Body.Repository.Name,
		*request.Body.Number,
		gitea.CreateIssueCommentOption{
			Body: "IssueInRedmine: https://projects.sdzecom.de/issues/" + fmt.Sprintf("%d", issue.ID),
		},
	)
	if err != nil {
		return issue.ID, -1, fmt.Errorf("error commenting issue with error: %w", err)
	}

	return issue.ID, comment.ID, nil
}

// NewEchoServer creates a new server with the given redmine project key, redmine context and gitea client
func NewEchoServer(rmpk string, rmCfg map[string]int64, rm *redmine.Context, gitea *gitea.Client) Server {
	e := echo.New()
	server := Server{
		Gitea:             gitea,
		Redmine:           rm,
		RedmineConfig:     rmCfg,
		RedmineProjectKey: rmpk,
		Echo:              e,
	}
	api.RegisterHandlers(e, api.NewStrictHandler(
		server,
		// add middlewares here if needed
		[]api.StrictMiddlewareFunc{},
	))

	return server
}
