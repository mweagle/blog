---
title: "Scaling the serverless summit requires environmental sympathy with dev & ops"
author: "Matt Weagle"
date: 2017-05-19T05:37:48.855Z
lastmod: 2023-02-18T21:57:32-08:00
tags: ["serverless", "sparta"]
description: ""
categories: 
    - Medium
subtitle: "The only way to gain confidence that a feature branch will work in the cloud is to run it in the cloud — with environmental sympathy"

image: "/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/1.png"
images:
 - "/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/1.png"
 - "/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/2.png"
 - "/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/3.png"
 - "/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/4.png"

---

#### The only way to gain confidence that a feature branch will work in the cloud is to run it in the cloud — with environmental sympathy

![image](/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/1.png#layoutTextWidth)


In the wake of Serverlessconf 2017 in Austin, there’s been an increasing number of discussions about today’s cold reality of serverless. While we can see the glory of serverless manifesting in the not-too distant future, the community still finds it [difficult to test](https://read.acloud.guru/testing-and-the-serverless-approach-495cef7495ea), deploy, debug, self-discover, and generally develop serverless applications.

The discussion has been amplified in recent days with [tweet storms](https://medium.com/@tobyhede/thoughts-on-the-serverless-development-experience-4023153757ee) summarized by [Toby Hede](https://medium.com/u/e14aa396c3c), and the great threads on the [Serverless Slack channel](http://wt-serverless-seattle.run.webtask.io/serverless-forum-signup?webtask_no_cache=1) from [Paul Johnston](https://medium.com/u/139b9858e816) that prompted this post. The common sentiment is that the difficultly with serverless gets more acute when developing applications composed of multiple sets of functions, infrastructure pieces, and identities evolving over time.

On the one hand, the [serverless approach to application architecture](https://read.acloud.guru/serverless-the-future-of-software-architecture-d4473ffed864) does implicitly address some of the high-availability aspects of service resiliency. For instance, you cloud assume — without empirical evidence — that AWS transparently migrates Lambda execution across Availability Zones in the face of localized outages. This is unlike a more traditional VM/container model, where you must explicitly distribute compute across isolated failure domains and load balance at a higher logical level (e.g. ELB and ALB).

While this intrinsic reliability is undoubtedly a good thing, overall resiliency isn’t so easily satisfied. Take for instance the canonical “Hello Serverless” application: an event-based thumbnailing workflow. Clients upload an image to an S3 bucket, a Lambda function handles the event, thumbnails the image, and posts it back to S3. Ship it.

Except, how do you actually test for the case when the [S3 bucket is unavailable](https://aws.amazon.com/message/41926/)? Or can you? I’m not thinking of testing against a _localhost_ mock API response, but the **actual** S3 bucket API calls — the bucket you’re accessing in production, via a dynamically injected environment variable.

Another example is when you have two Lambda functions, loosely coupled. The functions are blissfully ignorant of one another, although they share a mutual friend: Kinesis. In this use case, “Function A” publishes a message, perhaps with an embedded field whose value is another service’s event format (like an S3 triggering event) that’s consumed by “Function B”. While there’s no physical coupling, there’s potentially a deep logical coupling between them — one which might only appear at some future time as message contents drift across three agents in the pipeline.

How can we guard against this? How can we be certain about the set of functions which ultimately defines our service’s public contract?

_Are they coherent? Are the functions_ [_secure_](https://youtu.be/CiyUD_rI8D8?list=PLnwBrRU5CSTmruZzR8Z06j3pGglBZcdDr)_? Resilient? Correct? Observable? Scalable? How can we reduce uncertainty around_ [_non-functional requirements_](https://en.wikipedia.org/wiki/Non-functional_requirement)_?_

> [](https://twitter.com/mweagle/status/803070710578937856)


#### The non-functional requirements of serverless

The great thing about non-functional requirements is that they’re … _non-functional_. They speak to a system’s characteristics — how it should be — not what it should do, or how it should be done. In that sense, non-functional requirements both have nothing and everything to do with serverless.

![image](/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/2.png#layoutTextWidth)
The slide from Peter Bourgon’s [presentation on the microservice toolkit for Go](https://speakerdeck.com/peterbourgon/go-plus-microservices-equals-go-kit?slide=23)



The slide above is from Peter Bourgon’s excellent presentation on the design decisions behind [go-kit](https://github.com/go-kit/kit), a composable microservice toolkit for Go. The concerns listed apply equally to a JVM monolith, a Go-based set of microservices, or a NodeJS constellation supported by FaaS. If you’re running something in production, those *-_ilities_ lurk in the shadows whether or not they’re explicitly named.

In that sense, serverless is less a discontinuity with existing practice and more the next stage in the computing continuum — a theme emphasized in Tim Wagner’s [closing keynote](https://youtu.be/7ytNQTgqUXY?list=PLnwBrRU5CSTmruZzR8Z06j3pGglBZcdDr). It’s a technique that embeds more of the _*-ilities_ into the vendor platform itself, rather than requiring [secondary tools](https://netflix.github.io/). Serverless enables us to deliver [software faster](https://youtu.be/9tD5dX6hW0w?list=PLnwBrRU5CSTmruZzR8Z06j3pGglBZcdDr) and with fewer known unknowns — at least those that are externally observable.

Although serverless offloads more of these characteristics to the vendor, we still own the service. At the end of the day, each one of us is responsible to the customer, even when [conditions change](https://medium.com/@contact_16315/firebase-costs-increased-by-7-000-81dc0a27271d). We need to own it. And that means getting better at Ops. Or more specifically — _cloud-native development_.

![image](/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/3.png#layoutTextWidth)
Charity Majors does an excellent job describing the [operational best practices for serverless](https://charity.wtf/2016/05/31/operational-best-practices-serverless/)



#### The Base Camp — “Works on My Machine”

For many of us, the end result of our furious typing is in many cases a cloud-native application. In more mature organizations, our software constructs go through a structured [CI/CD](https://continuousdelivery.com/2014/02/visualizations-of-continuous-delivery/) pipeline and produce an artifact ready to ship. This artifact has a well-defined membrane through which only the purest configuration data flows and all dependencies are dynamic and well behaved.

On a day-to-day basis, though, there is often a lot of bash, docker-compose, DNS munging, and API mocks. There is also a lot of “works on my machine” — which may be true, at least at this instant — but probably doesn’t hold for everyone else on the team. And it definitely doesn’t provide a lot of confidence that it will work in the cloud.

The only way to gain confidence that a feature branch will work in the cloud is to run it in the cloud.
> Operations is the sum of all of the skills, knowledge and values that your company has built up around the practice of shipping and maintaining quality systems and software. — Charity Majors, [WTF is Serverless Operations](https://charity.wtf/2016/05/31/wtf-is-operations-serverless/)

If everyone on the team is developing their service feature branch in the cloud, complete with its infrastructure, then we’re all going to get better at ops. Because it’s development and ops rolled together. _And we’re all going to share a sense of Environmental Sympathy._

#### To the Summit — From #NoOps to #WereAllOps

Environmental Sympathy, inspired by [Mechanical Sympathy](https://www.infoq.com/presentations/mechanical-sympathy), is about applying awareness of our end goal of running in the cloud to the process of writing software.

![image](/posts/migrated/2017-05-19_scaling-the-serverless-summit-requires-environmental-sympathy-with-dev-ops/images/4.png#layoutTextWidth)


While it’s always been possible to provision isolated single-developer clusters complete with VMs, log aggregators, monitoring systems, feature flags, and the like, in practice it’s pretty challenging and expensive. And perhaps most aggravating, it can be very slow. Short development cycles are critical to developer productivity and that’s not really a hallmark of immutable, VM-based deploys.

_Serverless, precisely because it’s so heavily reliant on pre-existing vendor services and billed like a utility, makes it possible for every developer to exclusively develop their “service” in the cloud._

The service can have its own persistence engine, cache, queue, monitoring system, and all the other tools and namespaces needed to develop. Feature branches are the same as production branches and both are cloud-native by default. If during development, the _*-ilities_ tools prove too limiting, slow, or opaque, developer incentives and operational incentives are aligned. Together we build systems that make it easier to _ship and maintain quality systems and software_. Which will also help to minimize [MTTR](https://blog.heptio.com/perspective-on-multi-cloud-part-3-of-3-availability-and-multi-cloud-5018762d2702) as well.

Serverless, for both financial and infrastructure reasons, makes it possible to move towards cloud-native development and Environmental Sympathy. It represents a great opportunity to bring Dev and Ops — and QA, and SecOps) together. This allows us to mov from “worked on my machine” to “works in the cloud — I’ll slack you the URL.”

_From #NoOps to #WereAllOps._
