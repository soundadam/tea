---
title: tea
aliases:
  - /projects/teaway/
excerpt: 合上盖子也让 Mac 保持唤醒，用完把睡眠设置原样还回去。
weight: 10
presentation:
  category: macOS / CLI
  label: 用 Homebrew 安装
  command: brew install soundadam/tap/tea
  note: 没有守护进程，没有账号，没有遥测。
registry:
  github: soundadam/tea
  homebrew: soundadam/tap/tea
  docs: https://teaway.mintlify.app
release:
  version: "0.5.0"
  license: MIT
modules:
  - type: promo
    description: 让 Mac 保持唤醒——合着盖子、插不插电、不用开着终端都行——结束时只把自己改过的睡眠设置恢复原值。
    line: "{title} {version} 是 {license} 开源的 macOS 命令行客户端，支持 macOS 13 及以上，由 Homebrew 从带标签的源码构建。"
    actions:
      - label: 文档
        url: "{docs}"
      - label: GitHub
        url: "https://github.com/{github}"

  - type: snapshot
    title: 为什么不用 caffeinate
    body: macOS 自带 `caffeinate`，人坐在电脑前时它够用。但它只在自己的进程活着时阻止睡眠：关掉终端、SSH 断线，Mac 就可能睡过去；MacBook 合上盖子也照样睡。{title} 改的是系统的 `disablesleep` 设置，合盖照常运行，也不需要任何进程一直开着。改之前先记下原值，`off` 只恢复这个值，不动别的。
    note: 名字的由来：tea 是茶，对应 caffeinate 的咖啡因。0.4.2 及以前叫 teaway。
    metrics:
      - value: 合盖
        label: 照常唤醒，电池或电源都行
      - value: 还原
        label: "`off` 恢复它记下的原值"
      - value: "1"
        label: 个定时关机，可以精确取消
      - value: "{version}"
        label: "当前版本 · {license}"

  - type: section
    title: 把闲置的 Mac 当服务器用
    index: 1/2
    lede: 闲置 MacBook 当构建机、家里的 Mac mini 跑自托管服务，或者一次比你的会话更长的渲染、备份、传输。运行 `tea` 打开客户端，脚本里用 `tea on`、`tea shutdown after 2h`、`tea off`。

  - type: section
    title: 它不做什么
    index: 2/2
    lede: tea 只管电源状态，不配置远程登录、VPN、launchd 服务，也不管你跑的任务本身。合盖的 MacBook 散热更差：放在硬质、敞开的台面上，别放进包里。
---
