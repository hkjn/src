#!/usr/bin/env bash

set -euo pipefail

declare PASSWORD_MAIN_KEY=${PASSWORD_MAIN_KEY:-""}
declare PASSWORD_RECIPIENTS=${PASSWORD_RECIPIENTS:-""}

if [[ ! "${PASSWORD_MAIN_KEY}" ]]; then
	echo "No PASSWORD_MAIN_KEY specified." >&2
	exit 1
fi

echo "Importing public keys.."
for key in /etc/keys/*.asc;
	do gpg --import < $key
done

declare OWNER_TRUST=""
for PASSWORD in ${PASSWORD_RECIPIENTS}; do
	OWNER_TRUST="${OWNER_TRUST}${PASSWORD}:6:"$'\n'
done
OWNER_TRUST=${OWNER_TRUST%?}

echo "Importing owner trust settings for keys.."
echo "${OWNER_TRUST}" | gpg --import-ownertrust

echo "Importing private key.."
if [[ ! -f "/etc/secrets/keys/${PASSWORD_MAIN_KEY}.key" ]]; then
	echo "No such GPG key file: /etc/secrets/keys/${PASSWORD_MAIN_KEY}.key"
	exit 1
fi
gpg --import /etc/secrets/keys/${PASSWORD_MAIN_KEY}.key

echo "Dropping GPG identities from agent.."
echo RELOADAGENT | gpg-connect-agent

exec "$@"
