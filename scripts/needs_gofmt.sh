#!/bin/sh
#
# Checks if any files about to be committed need gofmt'ing.

# TODO: create go_git_hooks.sh, defining all functions needed.

echo "Checking if any files need gofmt.." >&2
IFS=$'\n'
if git rev-parse HEAD >/dev/null 2>&1; then
		FILES=$(git diff --cached --name-only | grep -e '\.go$');
else
		FILES=$(git ls-files -c | grep -e '\.go$');
fi
for file in $FILES; do
		echo "Checking $file.."
		badfile="$(gofmt -l "$file")"
		if test -n "$badfile" ; then
				echo "git pre-commit check failed: file needs gofmt: $file" >&2
				exit 1
		fi
done
