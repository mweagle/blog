---
title: "Mastodon - 2024-02-10T20:02:46Z"
subtitle: ""
canonical: https://hachyderm.io/users/mweagle/statuses/111908969970709060
description:
image: "/images/mastodon.png"

date: 2024-02-10T20:02:46Z
lastmod: 2024-03-09T15:24:03-08:00
image: ""
tags: ["Social Media"]

categories: ["mastodon"]
---
![Mastodon](/images/mastodon.png)

<p>Maybe this is my &quot;get off my lawn” moment, but sending an unstructured LLM text request that includes:</p><p>&quot;The output should be formatted as a JSON instance that conforms to the JSON schema below.”</p><p>followed by a markdown-escaped partial JSONSchema doesn&#39;t spark joy.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111908969970709060)

___

<p>API compatibility of `void* oracle(void*)` likely good news for the SRE world.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111909000390495306)

___

<p>I&#39;m using llama2 from <a href="https://ollama.com/" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://</span><span class="">ollama.com/</span><span class="invisible"></span></a> (to avoid the OpenAI r/t) and only changing the parameterized prompt. My (term, description) pairs are in YAML and serialized the same.</p><p>The response is either a JSONSchema envelope with the found results in the `properties.&lt;foo&gt;` element, or at the root in a `.&lt;foo&gt;` element.  Both responses include additional explanatory text. Submitting both in a prompt request produces another response format.</p><p>I&#39;m guessing this is expected</p><p>Going outside now</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/111909230349628707)

___
