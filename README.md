# stanford-rc/minio-console

**This is a modified version of the MinIO Console.** It is not MinIO, Inc.
software, it is not supported by MinIO, Inc., and it has not been reviewed or
endorsed by them.

Modified by The Board of Trustees of the Leland Stanford Junior University,
Stanford Research Computing, between **2026-04-09 and 2026-06-04**.

This program is free software, released under the **GNU Affero General Public
License version 3**, the same license as the works it derives from. See `LICENSE`
for the full text and `NOTICE` for attribution. Copyright in the unmodified
portions remains with MinIO, Inc., with Georg Mangold, and with the upstream
contributors.

MinIO(R) is a registered trademark of MinIO, Inc. It is used here only to
identify the software this fork is derived from.

## Provenance

This fork has two ancestors, and both matter for attribution:

1. `github.com/minio/object-browser`, the original MinIO Console, which MinIO,
   Inc. no longer develops in the open.
2. `github.com/georgmangold/console`, a fork of that work maintained by Georg
   Mangold. This repository is forked from it at release **v2.0.4**, upstream
   commit `2017f33b2`, dated 2025-09-05.

The Go module path was changed to `github.com/stanford-rc/minio-console` so that
this fork cannot be confused with either ancestor in a build.

## Why this fork exists

Stanford Research Computing serves this console to Elm users. The changes are
presentation-level, so that a user can tell which service they are logged into
and can reach the source for the modified server they are talking to.

## What was changed

- the login page carries a Stanford Research Computing banner
- the page header was reduced to a single line
- a link to the `elm-minio` source was added
- a stale orphaned screen was removed
- front-end assets rebuilt to match

`git log` is the authoritative list.

## Source availability

The complete corresponding source for this modified version is this repository.

## Upstream documentation

This README replaces the one inherited from `georgmangold/console`, which carried
that project's CI badges, release links and download counters. Those describe
that project and not this fork, which is why they are gone rather than edited.

For build and usage documentation, see the ancestors:

- https://github.com/georgmangold/console
- https://github.com/minio/object-browser
