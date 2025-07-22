---
title: "Mastodon - 2024-04-05T00:39:47Z"
subtitle: ""
canonical: https://hachyderm.io/users/mweagle/statuses/112215823989566555
description:
image: "/images/mastodon.png"

date: 2024-04-05T00:39:47Z
lastmod: 2024-04-05T00:39:47Z
image: ""
tags: ["Social Media"]

categories: ["mastodon"]
# generated: 2025-07-21T21:15:38-07:00
---
![Mastodon](/images/mastodon.png)

<p>Started reading <a href="https://www.cisa.gov/sites/default/files/2024-04/CSRB_Review_of_the_Summer_2023_MEO_Intrusion_Final_508c.pdf" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://www.</span><span class="ellipsis">cisa.gov/sites/default/files/2</span><span class="invisible">024-04/CSRB_Review_of_the_Summer_2023_MEO_Intrusion_Final_508c.pdf</span></a> and the ES states: &quot;The Board finds that this intrusion was preventable and should never have occurred.”</p><p>Is there ever an unpreventable intrusion that should have occurred?</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215823989566555)

___

<p>Sounds reasonable and I suspect is at odds with existing organizational &amp; personal incentives. </p><p>“The Board recommends that Microsoft’s CEO hold senior officers accountable for delivery against this plan. In the meantime, Microsoft leadership should consider directing internal Microsoft teams to deprioritize feature developments across the company’s cloud infrastructure and product suite until substantial security improvements have been made in order to preclude competition for resources.”</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215851528630539)

___

<p>Key rotation is a notoriously deferrable task, until it isn’t:</p><p>“Finally, this 2016 MSA key was originally intended to be retired in March 2021, but its removal was delayed due to unforeseen challenges associated with hardening the consumer key systems.”</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215865310101317)

___

<p>“Microsoft continued to rotate consumer MSA keys infrequently and manually until it stopped the rotation entirely in 2021 following a major cloud outage linked to the manual rotation process.”</p><p>Possibly security corollary to Lorin’s Law: <a href="https://surfingcomplexity.blog/2017/06/24/a-conjecture-on-why-reliable-systems-fail/" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://</span><span class="ellipsis">surfingcomplexity.blog/2017/06</span><span class="invisible">/24/a-conjecture-on-why-reliable-systems-fail/</span></a></p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215901981205904)

___

<p>There’s a theme that the enhanced logging tooling (G5 - TIL) used to produce evidence that confirmed the intrusion was both (a) an additional license that many don’t purchase and (b) a source of a tremendous amount of data that was challenging to decipher.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215935240512360)

___

<p>Electronic notifications sent to victims were ignored by some who “told FBI that they viewed these notifications as possible spam and disregarded them.”</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215953454325997)

___

<p>&quot;It requires a security-focused corporate culture of accountability, which starts with the CEO, to ensure that financial or other go-to -market factors do not undermine cybersecurity and the protection of Microsoft’s customers.”</p><p>Reminds me of Paul O’Neill’s kickoff meeting about workplace safety: <a href="https://davidburkus.com/2020/04/how-paul-oneill-fought-for-safety-at-alcoa/" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://</span><span class="ellipsis">davidburkus.com/2020/04/how-pa</span><span class="invisible">ul-oneill-fought-for-safety-at-alcoa/</span></a></p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112215993856177988)

___

<p>I would have liked to have been a fly on the wall in this architectural design review meeting. Shout out to everyone in that meeting who argued in favor of isolation ✊: </p><p>“Further, if Microsoft had not made the error that allowed consumer keys to authenticate to enterprise customer data”</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112216014609338555)

___

<p>Security as an optional &amp; paid upgrade doesn’t set customers up for success: “Security-related logging should be a core element of cloud offerings and CSPs should provide customers the foundational tools that provide them with the information necessary to detect, prevent, or quantify an intrusion”</p><p>I’ve found Colm’s hierarchy of design tradeoffs super helpful (apologies for the Melon link): <a href="https://x.com/colmmacc/status/986286693572493312" target="_blank" rel="nofollow noopener noreferrer" translate="no"><span class="invisible">https://</span><span class="ellipsis">x.com/colmmacc/status/98628669</span><span class="invisible">3572493312</span></a></p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112216035300856891)

___

<p>A lot of good suggestions for cell-based, strongly isolated, time-bound token design complemented by actual usage metrics that can be used to signal compromise.  </p><p>Premortem prompt: we’ve been compromised. How do we know? How fast can we know? How far in the past can we look?</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112216075331899081)

___

<p>One of the Ironies of Notifications is that people disregard notifications or view them as malicious: “In this intrusion, the Board found that some victims ignored or did not see the notifications, and some who saw them believed them to be spam or phishing”</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112216099448931426)

___

<p>There are a lot of good suggestions in this writeup and CISA has done a lot of work reconstructing the timeline. I have objections to hindsight counterfactuals (eg: “If Microsoft had not paused manual rotation of keys; if it had completed the migration of its MSA environment to rotate keys automatically; …”) but overall there are constructive next steps.</p><p>I strongly suspect there are people inside MSFT at the sharp end who were and are already aware of these issues.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112216131514839087)

___

<p>I hope this report provides them the support and space to make the changes they’ve been arguing for.</p>


###### [Mastodon Source 🐘](https://hachyderm.io/@mweagle/112216135775595699)

___
