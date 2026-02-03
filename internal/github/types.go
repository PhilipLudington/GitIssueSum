package github

import "time"

type Issue struct {
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	User        User          `json:"user"`
	Labels      []Label       `json:"labels"`
	Comments    int           `json:"comments"`
	CreatedAt   time.Time     `json:"created_at"`
	PullRequest *PullRequest  `json:"pull_request,omitempty"`
}

type User struct {
	Login string `json:"login"`
}

type Label struct {
	Name string `json:"name"`
}

type PullRequest struct {
	URL string `json:"url"`
}
