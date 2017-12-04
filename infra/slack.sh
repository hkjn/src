#!/bin/bash
#
# Slack library for bash scripts.
#

declare SLACK_TOKEN="${SLACK_TOKEN:-""}"
declare SLACK_URL="https://hooks.slack.com/services"

slacksend() {
	if [[ ! "${SLACK_TOKEN}" ]]; then
		echo "No SLACK_TOKEN available in environment." >&2
		return 1
	fi
	local message
	message="${1}"

	if [[ "$#" -ne 1 ]]; then
		echo "Usage is $0 'Message for slack'" >&2
		return 1
	fi

 	echo "Sending a message '${message}' to slack"
	local response
	response=$(curl -s -H 'Content-type: application/json' \
	                --data "{\"text\":\"$message\"}" \
			${SLACK_URL}/${SLACK_TOKEN})

    	if [[ "${response}" != "ok" ]]; then
	      echo "Bad Slack response: '${response}'" >&2
	      return 1
	fi
	return 0
}
