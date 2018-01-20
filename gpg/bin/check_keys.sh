#!/usr/bin/env bash
#
# Unlocks encrypted file given as first argument for viewing/editing.
#
# If the contents of the file changed, the cleartext file is re-encrypted.
#
# Regardless, the plaintext file is securely removed as the editor is
# closed, and is stored on tempfs only in the meanwhile.
#
# Unencrypted files can be encrypted with:
# cat holdings.json | unlock.sh holdings.json.gpg
#
declare BASE=/crypt
declare PASSWORD_SUB=${PASSWORD_SUB:-""}

set -euo pipefail

cd ${BASE}

source "/usr/local/bin/logging.sh"

cleanup() {
	info "Dropping GPG identities from agent.."
	echo RELOADAGENT | gpg-connect-agent
}

declare PASSWORD_RECIPIENTS=${PASSWORD_RECIPIENTS:-""}

if [[ "${PASSWORD_SUB}" ]]; then
	if [[ ! "${PASSWORD_RECIPIENTS}" ]]; then
		fatal "No PASSWORD_RECIPIENTS specified for subdirectory '${PASSWORD_SUB}'."
	fi
	debug "Using subdirectory ${PASSWORD_SUB} and recipients ${PASSWORD_RECIPIENTS}.."
fi

trap cleanup EXIT

debug "Using recipients ${PASSWORD_RECIPIENTS}"
cd ${PASSWORD_SUB}
for f in *.pgp; do
	info "${f} encrypted with:"
	set +e
	actual_recipients=$(gpg --batch --list-packets ${f} 2>&1 | grep 'encrypted with')
	set -e
	for recipient in ${PASSWORD_RECIPIENTS}; do
		if echo "${actual_recipients}" | grep -q ${recipient}; then
			info "${recipient}"
		else
			error "Missing recipient '${recipient}' for '${f}'"
		fi
	done
done
