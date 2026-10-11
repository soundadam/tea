---
title: "Release readiness"
description: "What a tea release publishes and what to check before tagging."
hidden: true
---

# Release readiness

## What a release is

Pushing a `v<version>` tag runs [`release.yml`](../.github/workflows/release.yml):

- it tests, then builds `tea-<version>-macos-aarch64.tar.gz` and
  `tea-<version>-macos-x86_64.tar.gz` with [`scripts/package.sh`](../scripts/package.sh)
  (ad-hoc signed executable, README, LICENSE, NOTICE, THIRD_PARTY_NOTICES);
- it checks that the packaged executable reports `tea <version>`;
- it creates the GitHub Release with the archives attached, notes taken from the
  `CHANGELOG.md` section for that version;
- for a tag without `-`, it writes [`packaging/homebrew/tea.rb`](../packaging/homebrew/tea.rb),
  filled with the version and both SHA-256 values, to `Casks/tea.rb` in
  `soundadam/homebrew-tap`. Tags with `-` are prereleases and leave the tap alone.

The version exists only in the tag. Naming, update behaviour, and the Cask
conventions are the [appkit](https://github.com/soundadam/appkit) contract.
Releases are immutable: a broken release is fixed by a new version.

## Before tagging

- `CHANGELOG.md` has a `## [<version>]` section (the workflow refuses a release
  tag without one).
- `main` CI is green, including the packaging smoke test.
- Real hardware verifies repeated `on`/`off`, external-state preservation,
  shutdown commit/status/cancel, and final restoration to `disablesleep=0`,
  with no `tea` shutdown event left.

## After the release

- `brew update && brew upgrade --cask soundadam/tap/tea` installs the new
  version, and `tea version` reports it.
- `tea auth register` refreshes the local helper; `tea auth status` reports it
  registered at the new version.
