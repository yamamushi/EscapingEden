package edenutil

import (
	"io/ioutil"
	"strings"
)

type BlackListTypeID int

const (
	BlackListIPs BlackListTypeID = iota
	BlackListUsernames
	BlackListEmails
)

func CheckBlacklist(input string, blacklist BlackListTypeID) bool {

	path := "assets/blacklist"

	switch blacklist {
	case BlackListIPs:
		path = path + "/" + "ips.txt"
	case BlackListUsernames:
		path = path + "/" + "usernames.txt"
	case BlackListEmails:
		path = path + "/" + "emails.txt"
	}

	// Open our blacklist file - hardcoded for now, I want to change this later
	blfile, err := ioutil.ReadFile(path)
	if err != nil {
		return false // If we can't open the file, we don't have a blacklist
	}
	// Split the file into lines
	blacklistLines := strings.Split(string(blfile), "\n")
	// Check if the email is in the blacklist
	for _, line := range blacklistLines {
		// if line begins with a # we ignore it
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.TrimSpace(line) == input {
			return true
		}
	}
	return false
}
