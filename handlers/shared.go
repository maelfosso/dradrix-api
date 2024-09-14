package handlers

import (
	"fmt"
	"strings"
)

func updateCurrentStatus(userCurrentStatus, justDoneStatus string) string {
	if userCurrentStatus == "is-creating-account" && justDoneStatus == "account-checked" {
		return "is-creating-account/set-profile"
	}
	if userCurrentStatus == "is-creating-account/set-profile" && justDoneStatus == "profile-set-up" {
		return "is-creating-account/set-org"
	}
	if userCurrentStatus == "is-creating-account/set-org" && justDoneStatus == "org-set-up" {
		return "registration-complete"
	}

	return userCurrentStatus
}

func nextLocation(path, currentStatus string) string {
	splittenPath := strings.Split(path, "/")

	endPath := ""
	switch currentStatus {
	case "is-creating-account/set-profile":
		endPath = "profile"
	case "is-creating-account/set-org":
		endPath = "org"
	case "registration-complete":
		return ""
	}

	if strings.HasPrefix(path, "/auth") {
		return fmt.Sprintf("/auth/%s", endPath)
	}
	if strings.HasPrefix(path, "/join") {
		return fmt.Sprintf("/join/%s/%s", splittenPath[2], endPath)
	}

	return ""
}
