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
declare KEEP_KEY_IN_MEMORY=${KEEP_KEY_IN_MEMORY:-""}

set -euo pipefail

cd ${BASE}
source "/usr/local/bin/logging.sh"

cleanup() {
	info "Securely removing '${CLEAR}'.."
	srm -v ${CLEAR}*
	if [[ ! "${KEEP_KEY_IN_MEMORY}" ]]; then
		info "Dropping GPG identities from agent.."
		echo RELOADAGENT | gpg-connect-agent
	fi
}

[[ "$#" -eq 1 ]] || fatal "Usage: $0 [encrypted file]"
declare TARGET=${1}
declare CRYPT=${BASE}/${TARGET}
declare PASSWORD_RECIPIENTS=${PASSWORD_RECIPIENTS:-""}
declare CLEAR=$(mktemp)

if [[ "${PASSWORD_SUB}" ]]; then
	CRYPT=${BASE}/${PASSWORD_SUB}/${TARGET}
	if [[ ! "${PASSWORD_RECIPIENTS}" ]]; then
		fatal "No PASSWORD_RECIPIENTS specified for subdirectory '${PASSWORD_SUB}'."
	fi
	debug "Using subdirectory ${PASSWORD_SUB} and recipients ${PASSWORD_RECIPIENTS}.."
fi

trap cleanup EXIT
[[ -e "$CRYPT" ]] || {
	info "No such file '$CRYPT', trying $CRYPT.pgp.."
	CRYPT="$CRYPT.pgp"
}

CHECKSUM_BEFORE=""
if [[ -e "$CRYPT" ]] && [[ ! -p /dev/stdin ]]; then
	info "Decrypting $CRYPT -> $CLEAR"
	debug "Cleartext file: $(ls -hsal ${CLEAR})"
	# TODO: may want to check that we can write to this file here, or the encryption
	# won't work.
	gpg --yes --output ${CLEAR} --decrypt /crypt/${PASSWORD_SUB}/$(basename ${CRYPT})
	if [[ $? -ne 0 ]]; then
		fatal "Error decrypting file."
	fi
	CHECKSUM_BEFORE=$(sha256sum $CLEAR)
else
	info "No such file '$CRYPT', creating new file '$CLEAR'"
fi

if [[ -p /dev/stdin ]]; then
	debug "/dev/stdin is a pipe, attempting to read it"
	cat > ${CLEAR}
else
	vim ${CLEAR}
fi

CHECKSUM_AFTER=$(sha256sum $CLEAR)
declare RECIPIENTS=""
for RECIPIENT in ${PASSWORD_RECIPIENTS}; do
	RECIPIENTS="${RECIPIENTS} --recipient ${RECIPIENT}"
done
debug "Using recipients ${RECIPIENTS}"

if [[ $CHECKSUM_BEFORE != $CHECKSUM_AFTER ]]; then
	info "Contents changed, re-encrypting ${CLEAR} -> $CRYPT"
	set +e
	$(gpg --yes --output /crypt/${PASSWORD_SUB}/$(basename ${CRYPT}) --encrypt --armor ${RECIPIENTS} ${CLEAR})
	STATUS=$?
	set -e
	if [[ ${STATUS} -ne 0 ]]; then
		fatal "Error encrypting file."
	fi
	debug "gpg --encrypt command exited with status '${STATUS}'"
fi

info "All done."
