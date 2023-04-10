---
title: "Lambda: Operations Delta"
author: "Matt Weagle"
date: 2016-03-30T15:26:18.905Z
lastmod: 2023-02-18T21:57:19-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "“That doesn’t look like what we saw yesterday.”"

image: "/posts/migrated/2016-03-30_lambda-operations-delta/images/1.jpeg"
images:
 - "/posts/migrated/2016-03-30_lambda-operations-delta/images/1.jpeg"

---

![image](/posts/migrated/2016-03-30_lambda-operations-delta/images/1.jpeg#layoutTextWidth)
Plateau Point Ahead



“That doesn’t look like what we saw yesterday.”

“There was a CVE announcement last night— what uses OpenSSL?”

“Seeing huge response time variance — WTF.”

“Is that node offline or just not receiving traffic?”

Once you start working on services, operations quickly become your ever-present companion. For as long as a service is online, operations are never “done”. Whether it’s performance, resilience, infrastructure changes, cost reduction, visibility, security or whatever else comes up, there’s always more to do. And it’s difficult to get those “known knowns” fixed when the team is burned out from responding to the “unknown unknowns” 3am Pager Duty alerts. [Paraphrasing Jeff Atwood](http://softwareengineeringdaily.com/2016/03/14/state-programming-jeff-atwood/), most of the problems aren’t (purely) code problems.

In the previous post I claimed that serverless marked a development discontinuity. It’s a major change to how services might be written, rather than a “lift-n-shift” development effort (or if the project isn’t going so well, “cram-n-slam”) to a shiny new OS/VM/package/container. Going serverless is about decomposing business value into a set of interlocking asynchronous, stateless event functions.

In many ways the pipes and filters model isn’t very novel, and I’m betting many of you have already thought about or even implemented something similar for your own services. What makes serverless different is that when the cloud vendor provides the platform, the architecture promises to offload some of the intrinsic operational costs. This benefit does come at a price, both financial and educational, but the potential gains are significant and should be considered in architectural discussions.

Serverless is a development model that acknowledges the critical importance of operations.

#### Reduce, Reuse, Redistribute
> You build it, you run it

— Werner Vogels, [ACM Queue,Volume 4, issue 4, June 30, 2006](http://queue.acm.org/detail.cfm?id=1142065)
> DevOps is not a job title; you cannot hire a “DevOp.” It is not a product; you cannot purchase “DevOps software.”

_— Limoncelli, Thomas A.; Chalup, Strata R.; Hogan, Christina J. The Practice of Cloud System Administration: Designing and Operating Large Distributed Systems, Volume 2 (p. 172). Pearson Education. Kindle Edition._

Keeping the packets running on time is a challenge even in the best of times. But when availability can literally be [compromised by sharks](https://www.youtube.com/watch?gl=SG&amp;v=1ex7uTQf4bQ&amp;hl=en-GB), we’re all going to need to work together. We can’t let the sharks win.

A collaborative, empathetic culture perspective intersects with the DevOps conversation. While I’m somewhat leery of the term, I completely agree with the notion that services must be designed and developed with operations in mind. How an organization facilitates that is itself an interesting topic, but for convenience let’s just the whole thing DevOps. It’s no longer acceptable to claim it worked in [dev, ops problem now](http://www.developermemes.com/2013/12/13/worked-fine-dev-ops-problem-now). That’s the primary reason why developers should run what they build. Internalizing development and operations costs is more likely to drive them down and bring about increased operational capability.

When development and operations are misaligned, business value flow is interrupted and one of your most important competitive advantages is lost: [speed](https://www.nginx.com/blog/adopting-microservices-at-netflix-lessons-for-team-and-process-design). It’s not possible to fail-fast/learn quickly if recovery is error prone and slow. Not to mention the early morning finger pointing meetings — those get old really fast.

In concert with the DevOps conversation, there’s also growing interest in the effectiveness of microservices to help business move faster. Decomposing a service into smaller functional units which can be handled by a small team improves overall agility. Smaller teams can iterate faster and at higher quality compared to larger ones, primarily due to significantly lower communication costs and related institutional overhead. There are fewer missed connections when everyone on the team can be in the same room sharing their two pizzas. This isn’t prescriptive, as small teams can be just as dysfunctional and ineffective as 20+ person groups delivering monoliths, but it certainly starts things off in a better place.

The microservice movement therefore promises to accentuate that fundamental competitive advantage: speed.

So we have more teams, producing more services, deploying more frequently, helping the business move faster. Wait, who’s on call again? Has anyone updated the alert trigger for database connection count? Which database in what environment &amp; geo? Forget it — I think we have a [heartbleed](http://heartbleed.com/).

We’ve just run into one of the most significant microservices trade-offs:
> You need a mature operations team to manage lots of services, which are being redeployed regularly.

— [Martin Fowler](http://martinfowler.com/articles/microservice-trade-offs.html)

Operations, in the traditional sense, typically aren’t fast and most definitely do not scale. Well, they do, but the transformation function is all wrong. Each service needs its own CI/CD pipeline, accounts, security review, testing environments/scripts/result parsers, supplementary shared cluster resources, vpcs, HA rollover/rollback, metrics, alerting, log aggregation, chat, escalation channels, etc. And in general these systems need to be _more_ available than the actual services they’re operating upon.

The only way this proliferation of microservices is going to work is if the mature operations team together with the developers (who potentially overlap with one another), backed by institutional support, starts to coalesce around some standard operating procedures. If every service is truly a bespoke one, operational costs will quickly come to dominate. It will take extreme persistence to prevent this conversation from stalling after infrastructure level recommendations, as trying to corral services back into the operational barn is [very difficult](https://twitter.com/frazelledazzell/status/713393137901633536).

But who says there’s even an operations team? Or even a single team — maybe the SecOps and the Chef/Ansible experts are in different groups. DevOps is about integrating operations and development, but with microservices we’ve vastly increased the operational costs and may not have a cost effective way to deal with it.

Enter serverless. Serverless is one approach to maintain business agility and speed of delivery. With reduced operational load, the business can continue to move fast, with high(er) confidence.

#### #LessOps
> No Server Is Easier To Manage Than No Server.

— [Werner Vogels AWS re:Invent 2015](https://www.youtube.com/watch?time_continue=1742&amp;v=y-0Wf2Zyi5Q)

Serverless doesn’t (and in fact, can’t) eliminate all operational costs, but it does promise to significantly reduce some of the more obvious ones. It’s about offloading some of the operational responsibility, not pretending there is some utopian cloud [#NoOps-land](http://cote.io/2016/03/24/sdt58/). Your service still needs [circuit-breakers](http://amzn.com/0978739213).

I’m most familiar with AWS Lambda, so this section is limited to their offerings, but I suspect that the competitor’s platforms offer something similar. I definitely need to learn more about Google’s platform — please shoot me a message [@](https://twitter.com/mweagle)mweagle if you have time to help me get up to speed.

A serverless application using AWS Lambda automatically includes:

*   OS updates/compliance: Implicit — AWS
*   HA/rollback: Implicit — AWS
*   Alerting: CloudWatch Alarms
*   Monitoring: CloudWatch, CloudWatch Logs + Lambda triggers
*   Log aggregation: CloudWatch Logs
*   Capacity Planning: on-demand/elastic (subject to AWS limits)
*   Security: granular IAM roles, API Gateway, Web Application Framework. Security Groups N/A.

There are competitors &amp; complementary tools for many of these dimensions, but the critical point is that there is something built-in, right from the start. Not to mention that the box is much smaller, since you’re not running your own EC2 instance. It doesn’t cover everything, and we do need new [serverless patterns](https://t.co/vLOn1Bjjo1), but it’s better than flying blind.

The AWS Lambda platform helps share some of your team’s operational responsibilities. In exchange for re-architecting your service in the Lambda model, you can significantly reduce your team’s non-customer facing workload.

In addition to these operational benefits, decomposing your service to a set of independently scheduled functions may even result in a significantly smaller codebase. Your code can focus more on the core business problem and less on whether the Ansible runtime is properly configured, whether you’re using all EC2 cores, or if the latest OpenSSL patch has been globally applied. Less code, less problems.

Serverless is about moving towards **#LessOps** — services that are operationally-aware by design and therefore require less of your team’s time to operate.

#### Plateau Ahead

Serverless marks a new plateau in service development. It may be difficult, if not ill-advised or even technically or financially impossible, for some services to consider a serverless architecture. Again, there are no silver bullets and serverless does come with its own downsides, some of which I’ll look at in another post. And needless to say, there’s ample evidence that operationally sustainable services are possible without a serverless architecture.

But for many services, serverless offers a huge reduction in operational costs and complexity together with an improvement in operational support. It’s still very early, and the plateau is a hazy image off on the horizon, but we’ve started the journey and will be making steady progress.
> [I]f you need to SSH into a server or an instance, you still have more to automate

— [Werner Vogels](http://www.allthingsdistributed.com/2016/03/10-lessons-from-10-years-of-aws.html)

Serverless is as much about eliminating SSH access and improving developer productivity as it is about reducing operational costs and increasing business agility. That is a powerful combination which promises significant rewards over the next few years — or more likely months given the rapid rate of innovation in the space.

### Next

I’m working to advance serverless adoption through [Sparta](http://gosparta.io), a Go framework for AWS Lambda services. Over the next series of posts I’ll walk through some code examples demonstrating how to build an AWS Lambda service using Sparta. I’ll step off the soapbox and actually get down to some code — stay tuned!_Photos by_ [_Willem van Bergen_](https://www.flickr.com/photos/willemvanbergen)
