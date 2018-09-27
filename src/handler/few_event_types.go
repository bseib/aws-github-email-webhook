// Cherry picked types from go-github/event_type.go so we can unmarshal the
// JSON for a github push event. Grabbed CommitAuthor from git_commits.go and
// User from users.go.
// 9/26/2018 bseib

// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// These event types are shared between the Events API and used as Webhook payloads.

package main

import (
	"time"
)

// CommitAuthor represents the author or committer of a commit. The commit
// author may not correspond to a GitHub User.
type CommitAuthor struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`

	// The following fields are only populated by Webhook events.
	Login *string `json:"username,omitempty"` // Renamed for go-github consistency.
}

func (c CommitAuthor) String() string {
	return Stringify(c)
}

// User represents a GitHub user.
type User struct {
	Login *string `json:"login,omitempty"`
	ID    *int64  `json:"id,omitempty"`
	// NodeID            *string    `json:"node_id,omitempty"`
	// AvatarURL         *string    `json:"avatar_url,omitempty"`
	// HTMLURL           *string    `json:"html_url,omitempty"`
	// GravatarID        *string    `json:"gravatar_id,omitempty"`
	Name *string `json:"name,omitempty"`
	// Company           *string    `json:"company,omitempty"`
	// Blog              *string    `json:"blog,omitempty"`
	// Location          *string    `json:"location,omitempty"`
	Email *string `json:"email,omitempty"`
}

func (u User) String() string {
	return Stringify(u)
}

// PushEvent represents a git push to a GitHub repository.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pushevent
type PushEvent struct {
	PushID       *int64            `json:"push_id,omitempty"`
	Head         *string           `json:"head,omitempty"`
	Ref          *string           `json:"ref,omitempty"`
	Size         *int              `json:"size,omitempty"`
	Commits      []PushEventCommit `json:"commits,omitempty"`
	Before       *string           `json:"before,omitempty"`
	DistinctSize *int              `json:"distinct_size,omitempty"`

	// The following fields are only populated by Webhook events.
	After      *string              `json:"after,omitempty"`
	Created    *bool                `json:"created,omitempty"`
	Deleted    *bool                `json:"deleted,omitempty"`
	Forced     *bool                `json:"forced,omitempty"`
	BaseRef    *string              `json:"base_ref,omitempty"`
	Compare    *string              `json:"compare,omitempty"`
	Repo       *PushEventRepository `json:"repository,omitempty"`
	HeadCommit *PushEventCommit     `json:"head_commit,omitempty"`
	Pusher     *User                `json:"pusher,omitempty"`
	Sender     *User                `json:"sender,omitempty"`
	//Installation *Installation        `json:"installation,omitempty"`
}

func (p PushEvent) String() string {
	return Stringify(p)
}

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	Message  *string       `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Distinct *bool         `json:"distinct,omitempty"`

	// The following fields are only populated by Events API.
	SHA *string `json:"sha,omitempty"`

	// The following fields are only populated by Webhook events.
	ID        *string       `json:"id,omitempty"`
	TreeID    *string       `json:"tree_id,omitempty"`
	Timestamp *Timestamp    `json:"timestamp,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
	Added     []string      `json:"added,omitempty"`
	Removed   []string      `json:"removed,omitempty"`
	Modified  []string      `json:"modified,omitempty"`
}

func (p PushEventCommit) String() string {
	return Stringify(p)
}

// PushEventRepository represents the repo object in a PushEvent payload.
type PushEventRepository struct {
	ID              *int64     `json:"id,omitempty"`
	NodeID          *string    `json:"node_id,omitempty"`
	Name            *string    `json:"name,omitempty"`
	FullName        *string    `json:"full_name,omitempty"`
	Owner           *User      `json:"owner,omitempty"`
	Private         *bool      `json:"private,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Fork            *bool      `json:"fork,omitempty"`
	CreatedAt       *Timestamp `json:"created_at,omitempty"`
	PushedAt        *Timestamp `json:"pushed_at,omitempty"`
	UpdatedAt       *Timestamp `json:"updated_at,omitempty"`
	Homepage        *string    `json:"homepage,omitempty"`
	Size            *int       `json:"size,omitempty"`
	StargazersCount *int       `json:"stargazers_count,omitempty"`
	WatchersCount   *int       `json:"watchers_count,omitempty"`
	Language        *string    `json:"language,omitempty"`
	HasIssues       *bool      `json:"has_issues,omitempty"`
	HasDownloads    *bool      `json:"has_downloads,omitempty"`
	HasWiki         *bool      `json:"has_wiki,omitempty"`
	HasPages        *bool      `json:"has_pages,omitempty"`
	ForksCount      *int       `json:"forks_count,omitempty"`
	OpenIssuesCount *int       `json:"open_issues_count,omitempty"`
	DefaultBranch   *string    `json:"default_branch,omitempty"`
	MasterBranch    *string    `json:"master_branch,omitempty"`
	Organization    *string    `json:"organization,omitempty"`
	URL             *string    `json:"url,omitempty"`
	ArchiveURL      *string    `json:"archive_url,omitempty"`
	HTMLURL         *string    `json:"html_url,omitempty"`
	StatusesURL     *string    `json:"statuses_url,omitempty"`
	GitURL          *string    `json:"git_url,omitempty"`
	SSHURL          *string    `json:"ssh_url,omitempty"`
	CloneURL        *string    `json:"clone_url,omitempty"`
	SVNURL          *string    `json:"svn_url,omitempty"`
}
