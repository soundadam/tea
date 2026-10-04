---
title: tea
aliases:
  - /projects/teaway/
excerpt: Keep a Mac awake with the lid closed, then put its sleep setting back.
weight: 10
presentation:
  category: macOS / CLI
  label: Install with Homebrew
  command: brew install soundadam/tap/tea
  note: No daemon. No account. No telemetry.
registry:
  github: soundadam/tea
  homebrew: soundadam/tap/tea
  docs: https://teaway.mintlify.app
release:
  version: "0.5.0"
  license: MIT
modules:
  - type: promo
    description: >-
      Keep a Mac awake — lid closed, on battery or AC, with no terminal left open — then restore exactly the sleep setting it changed.
    line: >-
      {title} {version} is an {license} command-line client for macOS 13 or later, built by Homebrew from the tagged source.
    actions:
      - label: Docs
        url: "{docs}"
      - label: GitHub
        url: "https://github.com/{github}"

  - type: snapshot
    title: Why not caffeinate
    body: >-
      macOS ships `caffeinate`, and it is the right tool while you sit at the Mac. It holds a sleep assertion only as long as its own process lives: close the terminal or drop the SSH session and the Mac may sleep. And closing a MacBook's lid still sleeps it. {title} switches the system's `disablesleep` setting instead, so the Mac stays up with the lid closed and nothing needs to keep running. Before changing it, {title} records the old value; `off` restores that value and nothing else.
    note: >-
      Why the name: tea, as opposed to caffeinate's caffeine. It was called teaway up to 0.4.2.
    metrics:
      - value: Lid closed
        label: Stays awake, battery or AC
      - value: Restore
        label: "`off` puts back the value it saved"
      - value: "1"
        label: Scheduled shutdown, cancellable exactly
      - value: "{version}"
        label: "Current release · {license}"

  - type: section
    title: Use a spare Mac as a server
    index: 1/2
    lede: >-
      A spare MacBook as a build box, a Mac mini at home, or a long render, backup or transfer that should outlive your session. Run `tea` for the client, or script it with `tea on`, `tea shutdown after 2h` and `tea off`.

  - type: section
    title: What it does not do
    index: 2/2
    lede: >-
      tea owns power state only. It does not set up Remote Login, a VPN, launchd services or the workload itself. A closed MacBook cools worse: keep it on a hard, open surface, never in a bag.
---
