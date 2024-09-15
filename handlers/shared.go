package handlers

import (
	"fmt"
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

func nextLocation(from, currentStatus string) string {
	endPath := ""
	switch currentStatus {
	case "is-creating-account/set-profile":
		endPath = "profile"
	case "is-creating-account/set-org":
		endPath = "org"
	case "registration-complete":
		return ""
	}

	return fmt.Sprintf("%s/%s", from, endPath)
}
