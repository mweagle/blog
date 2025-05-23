---
title: "Mastodon - 2024-02-12T06:26:05Z"
subtitle: ""
canonical: https://hachyderm.io/users/mweagle/statuses/111917083212359082
description:
image: "/images/mastodon.png"

date: 2024-02-12T06:26:05Z
lastmod: 2024-02-12T06:26:05Z
image: ""
tags: ["Social Media"]

categories: ["mastodon"]
# generated: 2025-05-22T22:29:20-07:00
---
![Mastodon](/images/mastodon.png)

<p>&quot;From the beginning I made some core decisions that the company has had to stick to, for better or worse, these past four years. This post will list some of the major decisions made and if I endorse them for your startup, or if I regret them and advise you to pick something else.&quot;</p><p><a href="https://cep.dev/posts/every-infrastructure-decision-i-endorse-or-regret-after-4-years-running-infrastructure-at-a-startup/" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://</span><span class="ellipsis">cep.dev/posts/every-infrastruc</span><span class="invisible">ture-decision-i-endorse-or-regret-after-4-years-running-infrastructure-at-a-startup/</span></a></p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111917083212359082)

___

<p>The author calls out a pain point I&#39;ve encountered and also invested in addressing in multiple teams:</p><p>&quot;With GitOps we’ve had to invest in tooling to help people answer questions like “I did a commit: why isn’t it deployed yet”.&quot;</p><p>- Product always wants to know when a feature will be available<br />- Dev always wants to know what&#39;s running when a test a test fails<br />- Ops always wants to know what version of code/artifact is running during an incident.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111917106177758271)

___

<p>- Security always wants to know if there is CVE risk</p><p>There are many stakeholders who need to reverse lookup from an `instance-id`, not only a `class`, yet so much context is lost during CI/CD. </p><p>Visualizing the code lineage multiverse is a gnarly problem beyond just the number of transformations &amp; packaging formats.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111917160333503715)

___

<p>At a mundane level, we will only know when things are better when we are able to signify the things. Can you step in the same river twice and whatnot.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111917179990847754)

___
