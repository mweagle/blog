---
title: "Serverless: The next discontinuity"
author: "Matt Weagle"
date: 2016-03-21T22:06:55.585Z
lastmod: 2023-02-18T21:57:17-08:00
tags: ["serverless", "sparta"]
ShowToc: true
description: ""

subtitle: "Over the past year or so, the “serverless” architecture movement has seen a frenzy of activity. AWS announced Lambda at re:Invent 2014 and…"

image: "/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/1.jpeg"
images:
 - "/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/1.jpeg"
 - "/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/2.png"
 - "/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/3.jpeg"
 - "/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/4.jpeg"

---


![image](/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/1.jpeg#layoutTextWidth)
Mabry Mill



Over the past year or so, the “serverless” architecture movement has seen a frenzy of activity. AWS announced Lambda at re:Invent 2014 and after a predictable period of cautious exploration by the developer community, began to really take off.

![image](/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/2.png#layoutTextWidth)
Google Trends for “serverless” — units not included



In addition to the continued improvements to [AWS Lambda](https://aws.amazon.com/releasenotes/AWS-Lambda), other signs of larger industry growth include:

*   Google’s recent announcement of [Google Cloud Functions](https://cloud.google.com/functions/docs)
*   Multiple open source efforts to ease the operational aspects of AWS Lambda including [serverless](https://github.com/serverless/serverless), [Apex](https://github.com/apex/apex), [Kappa](https://github.com/garnaat/kappa), and my own [Sparta](http://gosparta.io)
*   A serverless [conference](http://serverlessconf.io/)
*   A dedicated serverless [blog](https://serverlesscode.com/), AWS Lambda [book](http://amzn.com/B016JOMAEE), and another [Learn Serverless](http://justserverless.com/blog/releasing-our-learn-serverless-book/) book in the works.
*   Many Lambda applications, including [VoiceOps](https://github.com/1Strategy/alexa-aws-administration), [Tumblr-style blogs](https://github.com/matteobrusa/Tumbless), and [video transcoders](https://github.com/kefabean/lambda-transcoder).

For many types of applications, “serverless” marks a discontinuity in both application design and locus of control. It is a higher level of abstraction (eg, [_Your server as a function_](https://engineering.twitter.com/research/publication/your-server-as-a-function)) that enables developers to focus on their **application**, not the [undifferentiated muck](https://aws.amazon.com/blogs/aws/we_build_muck_s/) required to support it.

#### #WaterOps

I recently took a trip through the [Blue Ridge Parkway](http://www.blueridgeparkway.org/) in southern Virginia. One of the most visited sites along the drive is the famous [Mabry Mill](https://en.wikipedia.org/wiki/Mabry_Mill), a watermill, workshop, and blacksmith shop alongside a small stream. The mill is still in working order and it’s quite impressive to see the gear and band system transferring energy from the waterwheel to the jig saw, lathe, and other equipment.

![image](/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/3.jpeg#layoutTextWidth)


These machines were magnifiers for Mabry’s skills: everything needed for the farm could be made with these tools (well, he needed the blacksmith shop as well).

However, the machines needed power, and with unpredictable rainfall amounts, Mabry needed increasing amounts of land on which he could build an labyrinthine flume system to gather water and channel it to his waterwheel.

![image](/posts/migrated/2016-03-21_serverless-the-next-discontinuity/images/4.jpeg#layoutTextWidth)


These flumes didn’t deliver [logfiles](https://flume.apache.org/), but water to the wheel. For Mabry, these flumes were the _undifferentiated muck_ of his time. If a flume was leaking, or a gate was closed, or even if it hadn’t rained in a while, then the “saw service” was partially available or even offline. You might say that Mabry was an early and unwitting practitioner of #WaterOps. While I appreciated Mabry’s industriousness, I couldn’t help think that his type of mill was likely one of the last self-powered, as the [power grid](https://en.wikipedia.org/wiki/Electrical_grid) was starting to come into existence.

How much easier would it have been to plug his machines into a socket and have consistent, highly available power, without the need to replace a rotting board?

#### Functional Power

This is the (metaphorical) promise of serverless architectures: consistent, highly available compute power that enables an application to securely and reliably deliver business value. The freedom to focus your efforts on improving your application, without the requirement to replace a leaky SSH Security Group or an offline EC2 instance.

Serverless isn’t a silver bullet, and not all types of services could or even [should] (http://engineers.coffee/episodes/2016/2/)consider moving to this architecture. However, for many others the possibility of consistent, highly-available compute power offers a compelling alternative.

#### Future

I think the future of serverless architectures is bright, and am actively working on [Sparta](http://gosparta.io/)— a framework for Go-based AWS Lambda applications. Over the next series of posts I’ll dive more into the AWS Lambda ecosystem and how it opens up a new way to securely, reliably, and consistently power your applications.
