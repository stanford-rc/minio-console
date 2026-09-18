#!/usr/bin/env bash
#
# verify-standalone.sh builds, vets and unit-tests this repository on its own.
#
# WHY THIS NEEDS A WORKAROUND
#
# go.mod declares "module github.com/stanford-rc/minio-console", while every
# file in the tree imports "github.com/minio/console/...". Standalone that does
# not resolve, so a plain "go build ./..." fails with a wall of
#
#   no required module provides package github.com/minio/console/api/operations
#
# which looks like a broken checkout and is not.
#
# The mismatch is load bearing, not an oversight. github.com/stanford-rc/minio
# consumes this repository through a version replace:
#
#   replace github.com/minio/console => github.com/stanford-rc/minio-console v1.9.7-stanford-rc
#
# Go verifies that the fetched module declares the path it was replaced WITH, so
# go.mod has to say stanford-rc/minio-console. Meanwhile an import of
# github.com/minio/console/api resolves through the replace to this tree, so the
# tree's own imports have to keep saying minio/console. Neither side can change
# alone:
#
#   * rewrite the imports to stanford-rc/minio-console and the consumer breaks,
#     because that path is in its build list only as a replacement TARGET and
#     never as a required module, so the internal imports go unresolvable.
#   * change go.mod to say minio/console and the consumer breaks earlier still,
#     at fetch time, on the declared-path check.
#
# Both were tried. So rather than pick a side, this script aligns the module
# path to the imports for the duration of the checks and puts it back
# afterwards, including on failure or interrupt. Nothing is committed and the
# consumer is unaffected.
#
# Integration suites are excluded: sso-integration and integration expect a
# live server on localhost and fail here for reasons that have nothing to do
# with the code under test.

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

declared_path="github.com/stanford-rc/minio-console"
import_path="github.com/minio/console"

current="$(go list -m 2>/dev/null || true)"

restore_needed=0
backup=""

cleanup() {
	local rc=$?
	if [[ "${restore_needed}" -eq 1 ]]; then
		if cp -- "${backup}" go.mod; then
			echo "restored go.mod"
		else
			# Loud, because leaving the module path rewritten would make the
			# next consumer build fail in a way that points nowhere near here.
			echo "FAILED to restore go.mod from ${backup}; restore it by hand" 1>&2
			rc=1
		fi
	fi
	[[ -n "${backup}" ]] && rm -f -- "${backup}"
	exit "${rc}"
}

if [[ "${current}" == "${declared_path}" ]]; then
	backup="$(mktemp go.mod.verify.XXXXXX)"
	cp -- go.mod "${backup}"
	restore_needed=1
	trap cleanup EXIT INT TERM

	echo "temporarily aligning the module path to ${import_path} for verification"
	# Only the module line, and only if it is the first line, so a replace or
	# require naming the same path is left alone.
	sed -i "1s|^module ${declared_path}\$|module ${import_path}|" go.mod

	if [[ "$(go list -m)" != "${import_path}" ]]; then
		echo "module path was not aligned; go.mod line 1 is not what was expected" 1>&2
		exit 1
	fi
elif [[ "${current}" == "${import_path}" ]]; then
	# Already aligned, so somebody is mid-experiment or a previous run did not
	# restore. Verify, but do not touch go.mod and do not hide the state.
	echo "NOTE: go.mod already declares ${import_path}; verifying without changing it"
	echo "NOTE: it must say ${declared_path} when committed, or the consumer build breaks"
else
	echo "unexpected module path ${current:-<unknown>}; expected ${declared_path}" 1>&2
	exit 1
fi

# Build over ./... rather than an explicit list: "go build" rejects a package
# holding only test files, which ./... tolerates and this tree contains
# (replication). node_modules ships Go files of its own that ./... picks up,
# flatted/golang among them; they build cleanly and are left alone here rather
# than filtered, so that a genuine breakage in them stays visible.
echo "==> go build ./..."
go build ./...

# vet and test do take a list, to drop the integration suites, which expect a
# live server on localhost, and node_modules, which is not ours to report on.
pkgs="$(go list ./... | grep -vE '/(sso-)?integration$|/node_modules/' || true)"
if [[ -z "${pkgs}" ]]; then
	echo "no packages found to verify" 1>&2
	exit 1
fi

echo "==> go vet"
# shellcheck disable=SC2086
go vet ${pkgs}

echo "==> go test -short"
# shellcheck disable=SC2086
go test -short ${pkgs}

echo "==> verified"
