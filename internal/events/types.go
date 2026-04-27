package events

import "github.com/wanderer-npm/dispatch/internal/discord"

type user struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type repo struct {
	FullName    string `json:"full_name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type commitAuthor struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type commit struct {
	ID      string       `json:"id"`
	Message string       `json:"message"`
	URL     string       `json:"url"`
	Author  commitAuthor `json:"author"`
}

func authorOf(u user) *discord.Author {
	if u.Login == "" {
		return nil
	}
	return &discord.Author{
		Name:    u.Login,
		URL:     u.HTMLURL,
		IconURL: u.AvatarURL,
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	i := max - 1
	for i > 0 && s[i]&0xC0 == 0x80 {
		i--
	}
	return s[:i] + "…"
}

func repoLink(r repo) string {
	if r.HTMLURL != "" {
		return "[" + r.FullName + "](" + r.HTMLURL + ")"
	}
	return r.FullName
}
