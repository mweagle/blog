---
title: "Mastodon - 2024-02-27T01:22:33Z"
subtitle: ""
canonical: https://hachyderm.io/users/mweagle/statuses/112000824344693164
description:
image: "/images/mastodon.png"

date: 2024-02-27T01:22:33Z
lastmod: 2024-02-27T01:22:33Z
image: ""
tags: ["Social Media"]

categories: ["mastodon"]
# generated: 2024-10-23T18:04:53-07:00
---
![Mastodon](/images/mastodon.png)

<p>Been using LangChain and Chat GPT to analyze a corpus of text (<a href="https://info.arxiv.org/help/bulk_data.html" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://</span><span class="ellipsis">info.arxiv.org/help/bulk_data.</span><span class="invisible">html</span></a>). I&#39;m taking PDFs, extracting the plaintext per-page and querying for something like &quot;No more than N sentences that reference concept X&quot;. </p><p>It *seems* to do ok with general sentiment and fails to extract sentences that are related to the concept. It happily returns sentences that a human reader would consider unrelated.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112000824344693164)

___

<p>To &quot;debug&quot;, I try the interactive chat version to see if I can understand why certain sentences are returned.</p><p>It initially respected limiting responses to the input text. It acknowledged its error (?) and understandably couldn&#39;t explain why a specific sentence was included.</p>

![](374058e00762c5da.png)

###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112000856385658792)

___

<p>The refined sentences were not accurate, so I followed up with a request to limit responses to the input text only. </p><p>It then started to include sentences that did not exist in the text at all.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112000864225973862)

___

<p>Giving up now - maybe there&#39;s a &quot;magic prompt&quot; I should be using, but I have no clue how to find it.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112000916952389049)

___
