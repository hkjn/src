run_gcloud() {
	docker run --rm -it \
	           -e GCLOUD_PROJECT=zero-iosdev \
	           -e GCLOUD_SECRETS_PATH=/etc/secrets/default-zero-iosdev-editor.json \
	           -v $(pwd)/.gcp:/etc/secrets:ro \
	           hkjn/gcloud:1.0.0
}
alias gcd='run_gcloud $@'
