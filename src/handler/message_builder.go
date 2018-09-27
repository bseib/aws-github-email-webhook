package main

import "strings"

func BuildEmailSubject(event PushEvent) string {
	var buf strings.Builder
	buf.WriteString("[")
	buf.WriteString(*event.Repo.FullName)
	buf.WriteString("] ")
	if nil != event.Commits && len(event.Commits) > 0 {
		firstCommit := event.Commits[0]
		id := *firstCommit.ID
		buf.WriteString(string(id[0:6]))
		if hasMessage(firstCommit) {
			buf.WriteString(": ")
			msg := *firstCommit.Message
			index := strings.Index(msg, "\n")
			if index < 0 {
				buf.WriteString(msg)
			} else {
				buf.WriteString(msg[0:index])
			}
		}
	}
	return buf.String()
}

func BuildEmailMessageBody(event PushEvent) string {
	var buf strings.Builder
	buf.WriteString("  Branch: " + *event.Ref + "\n")
	buf.WriteString("  Home:   " + *event.Repo.URL + "\n")
	if nil != event.Commits {
		for _, commit := range event.Commits {
			buf.WriteString("  Commit: " + *commit.ID + "\n")
			buf.WriteString("      " + *commit.URL + "\n")
			buf.WriteString("  Author: " + *commit.Author.Name + " <" + *commit.Author.Email + ">\n")
			buf.WriteString("  Date:   " + commit.Timestamp.String() + "\n")
			buf.WriteString("\n")

			if hasChanges(commit) {
				buf.WriteString("  Changed paths:\n")
			}
			if hasModified(commit) {
				appendFileList(&buf, "M", commit.Modified)
			}
			if hasAdded(commit) {
				appendFileList(&buf, "A", commit.Added)
			}
			if hasRemoved(commit) {
				appendFileList(&buf, "R", commit.Removed)
			}
			if hasChanges(commit) {
				buf.WriteString("\n")
			}

			if hasMessage(commit) {
				buf.WriteString("  Log Message:\n")
				buf.WriteString("  -----------\n")
				appendMessage(&buf, commit.Message)
				buf.WriteString("\n")
			}

			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func hasChanges(commit PushEventCommit) bool {
	return hasModified(commit) || hasAdded(commit) || hasRemoved(commit)
}

func hasModified(commit PushEventCommit) bool {
	return nil != commit.Modified && len(commit.Modified) > 0
}

func hasAdded(commit PushEventCommit) bool {
	return nil != commit.Added && len(commit.Added) > 0
}

func hasRemoved(commit PushEventCommit) bool {
	return nil != commit.Removed && len(commit.Removed) > 0
}

func hasMessage(commit PushEventCommit) bool {
	return nil != commit.Message && len(*commit.Message) > 0
}

func appendFileList(buf *strings.Builder, prefix string, list []string) {
	for _, filename := range list {
		buf.WriteString("    ")
		buf.WriteString(prefix)
		buf.WriteString(" ")
		buf.WriteString(filename)
		buf.WriteString("\n")
	}
}

func appendMessage(buf *strings.Builder, message *string) {
	buf.WriteString("  ")
	buf.WriteString(*message)
	buf.WriteString("\n")
}
