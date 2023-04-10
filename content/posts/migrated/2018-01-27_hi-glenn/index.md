---
title: "Hi Glenn,"
author: "Matt Weagle"
date: 2018-01-27T04:11:21.600Z
lastmod: 2023-02-18T21:57:39-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "That error typically means that the go get -u -v ./... command didn’t properly fetch all the dependencies. I’ll take a look and see if I…"

---

Hi Glenn,

That error typically means that the `go get -u -v ./...` command didn’t properly fetch all the dependencies. I’ll take a look and see if I can repro. Glad you were able to work through it and please LMK if you run into any new issues.
