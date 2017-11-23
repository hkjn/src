# Golang git hooks, useful in constructing e.g. pre-commit, pre-push.

has_conflicts() {
  echo "Checking for merge conflicts.." >&2
  conflicts=$(git diff --cached --name-only -S'<<<<<< HEAD')
  if [[ -n "$conflicts" ]]; then
    echo "Unresolved merge conflicts in this commit:" >&2
    echo $conflicts >&2
    return 1
  fi
  return 0
}

needs_gofmt() {
  echo "Checking if any files need gofmt.." >&2
  IFS=$'\n'
  local files=""
  if [ git rev-parse HEAD >/dev/null 2>&1 ]; then
    files=$(git diff --cached --name-only | grep -e '\.go$');
  else
    files=$(git ls-files -c | grep -e '\.go$');
  fi
  local failed=0
  for file in $files; do
    if [[ ! -e "$file" ]]; then
      # If a cached file doesn't exist it must be about to be deleted
      # in this commit, so no need to check it.
      continue
    fi
    local badfile="$(gofmt -l "$file")"
    if test -n "$badfile" ; then
      echo "git pre-commit check failed: file needs gofmt: $file" >&2
      failed=1
    fi
  done
  if [[ $failed -ne 0 ]]; then
    return 1
  fi
  return 0
}

prevent_dirty_tree() {
  if [ "$#" -ne 2 ]; then
    echo "Usage: prevent_dirty_tree [directory] [comment]" >&2
    return 1
  fi
  if [ $(git status --porcelain 2>/dev/null ${1} | grep "^ M" | wc -l) -ne "0" ]; then
    echo "Diff in ${1} - ${2}:" >&2
    echo $(git diff --numstat ${1})
    return 1
  fi
  return 0
}

prevent_hacks() {
  echo "Checking for strings indicating hacks.." >&2
  local files=""
  if [ git rev-parse HEAD >/dev/null 2>&1 ]; then
    # TODO(hkjn): This also can return files that are removed in the
    # working tree, which we should not be trying to grep through..
    files=$(git diff --cached --name-only | grep -v vendor/)
  else
    files=$(git ls-files -c | grep -v vendor/)
  fi

  if [[ ! "$files" ]]; then
    return 0 # Nothing to check.
  fi
  local failed=0

  local baseDir="$(cd "$(dirname "$0")" && pwd)/../.."
  cd $baseDir

  for file in $files; do
    if grep -ni FIXME "$file"; then
      echo "^^ Please remove offending string 'FIXME' from '$file'." >&2
      failed=1
    fi
    if grep -ni "DO NOT SUBMIT" "$file"; then
      echo "^^ Please remove offending string 'DO NOT SUBMIT' from '$file'." >&2
      failed=1
    fi
  done

  if [[ $failed -ne 0 ]]; then
    return 1
  fi
  return 0
}

run_go_tests() {
  echo "Running all Go tests.." >&2
  targets=$(go list ./... 2>/dev/null | grep -v /vendor/)
  if [[ ! "$targets" ]]; then
    return 0 # Nothing to test
  fi

  output=$(go test -race $targets 2>&1)
  if [ $? -eq 0 ]; then
    return 0
  fi
  if echo "$output" | grep "matched no packages" >/dev/null; then
    # Special case for "there's no packages in this repo", which is fine.
    return 0
  fi
  echo "Go tests failed:" >&2
  echo "$output" >&2
  return 1
}

run_go_vet() {
  echo "Running Go vet command.." >&2
  local targets=$(go list ./... 2>/dev/null | grep -v /vendor/)
  if [[ ! "$targets" ]]; then
    return 0 # Nothing to test
  fi
  local output
  output=$(go vet $targets 2>&1)
  if [ $? -eq 0 ]; then
    return 0
  fi
  if echo "$output" | grep "matched no packages" >/dev/null; then
    # Special case for "there's no packages in this repo", which is fine.
    return 0
  fi
  echo "Go vet failed:\n$output" >&2
  return 1
}

update_bindata() {
  if [ -d "bindata/" ]; then
    echo "Checking if bindata needs regenerating.." >&2
    go-bindata -pkg="bindata" -o bindata/bin.go tmpl/
    # Note: since go-bindata doesn't properly format its output, we need
    # to do it ourselves.
    gofmt -w bindata/bin.go
    local dirty=$(prevent_dirty_tree bindata/ "commit bindata changes first")
    return ${dirty}
  fi
  return 0
}

update_godep() {
  if [[ -d "Godeps/" ]]; then
    echo "Checking for godeps updates.." >&2
    godep update ...
    local dirty=$(prevent_dirty_tree Godeps/ "commit dependency changes first")
    return ${dirty}
  fi
  return 0
}
