declare GCLOUD_PROJECT=${GCLOUD_PROJECT:-""}
declare GCLOUD_SECRETS_PATH=${GCLOUD_SECRETS_PATH:-""}

source "${HOME}/google-cloud-sdk/completion.bash.inc"
source "${HOME}/google-cloud-sdk/path.bash.inc"

fatal() {
	echo "[ .bash_profile ] $@" >&2
	exit 1
}

info() {
	echo "[ .bash_profile ] $@" >&2
}

gcloud_login() {
	[[ "${GCLOUD_PROJECT}" ]] || fatal "No GCLOUD_PROJECT specified"
	[[ "${GCLOUD_SECRETS_PATH}" ]] || fatal "No GCLOUD_SECRETS_PATH specified"
	[[ -e "${GCLOUD_SECRETS_PATH}" ]] || fatal "Specified GCLOUD_SECRETS_PATH '${GCLOUD_SECRETS_PATH}' doesn't exist"

	info "Setting project to '${GCLOUD_PROJECT}'.."
	gcloud config set project ${GCLOUD_PROJECT}
	info "Activating service account using '${GCLOUD_SECRETS_PATH}'.."
	gcloud auth activate-service-account --key-file="${GCLOUD_SECRETS_PATH}"
}

gcloud_login
