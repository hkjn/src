show_ignite_diff() {
	local target
	target=${1}
	if [[ ! -e bootstrap/${target}.json ]]; then
		echo "Unknown target ${target}." >&2
		return 1
	fi
	run_tf output hkjn_ignite_json_${target} > /tmp/${target}_output.json
	if [[ $? -ne 0 ]]; then
		echo "tf output command failed: $(cat /tmp/${target}_output.json)" >&2
	fi
	diff <(jq '.' < /tmp/${target}_output.json) <(jq '.' < bootstrap/${target}.json)
}

generate_ignite_configs() {
	echo "Running generate_ignite_configs.go tool to generate Ignite .json.."
	docker run --rm -it \
	           -v /etc/secrets/secretservice:/etc/secrets/secretservice:ro \
	           -v $(pwd):/home/go/src/hkjn.me/src/infra \
	           -w /home/go/src/hkjn.me/src/infra \
	       hkjn/golang go run generate_ignite_configs.go
}

fetch_checksums() {
	docker run --rm -it \
	           -v /etc/secrets/secretservice:/etc/secrets/secretservice:ro \
	           -v $(pwd):/home/go/src/hkjn.me/src/infra \
	           -w /home/go/src/hkjn.me/src/infra \
	       hkjn/golang go run fetch_checksums.go
}

run_tf() {
	local action
	action=$1
	if [[ "${action}" = plan ]]; then
		if [[ ! -e /etc/secrets/secretservice/seed ]]; then
			echo "FATAL: Missing /etc/secrets/secretservice/seed." >&2
			return 1
		fi
		if [[ ! -e /etc/secrets/secretservice/salt ]]; then
			echo "FATAL: Missing /etc/secrets/secretservice/salt." >&2
			return 1
		fi
		local sshash
		sshash=$(echo $(cat /etc/secrets/secretservice/seed)'|'$(cat /etc/secrets/secretservice/salt) | sha512sum | cut -d ' ' -f1)
		SECRETSERVICE_HASH=${sshash} generate_ignite_configs
		if [[ $? -ne 0 ]]; then
			echo "FATAL: ignite.go failed." >&2
			return 1
		fi
	fi
	echo "Running 'terraform $@'.."
	# TODO: Below we take current VERSION file, but could run an older version
	# for some targets, as specified in ignition.py..
	echo "version = \"$(cat VERSION)\"" > terraform.tfvars
	# TODO: By doing 'tf plan -detailed-exitcode', we can check for status 2 (there was
	# a diff), and if so, run the ignite_diff command above for a not-horrible comparison
	# on where the metadata differs.
	docker run --rm -it -v $(pwd):/home/tfuser \
	           -e GOOGLE_APPLICATION_CREDENTIALS=/home/tfuser/.gcp/tf-dns-editor.json \
		   hkjn/terraform $@
}

alias ignite_diff='show_ignite_diff $@'
alias tf='run_tf $@'
