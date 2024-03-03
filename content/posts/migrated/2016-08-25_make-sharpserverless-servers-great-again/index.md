---
title: "Make #serverless Servers Great Again"
author: "Matt Weagle"
date: 2016-08-25T05:17:06.697Z
lastmod: 2023-02-18T21:57:23-08:00
tags: ["serverless", "sparta"]
description: ""
categories: 
    - Medium
subtitle: "The serverless pay-as-you go, on-demand compute model promises the opportunity to focus more of your time on delivering business value and…"

image: "/posts/migrated/2016-08-25_make-sharpserverless-servers-great-again/images/1.png"
images:
 - "/posts/migrated/2016-08-25_make-sharpserverless-servers-great-again/images/1.png"
 - "/posts/migrated/2016-08-25_make-sharpserverless-servers-great-again/images/2.png"


---

The serverless pay-as-you go, on-demand compute model promises the opportunity to focus more of your time on delivering business value and less time managing infrastructure.

However, there are times when a more traditional application server model is appropriate. For instance,

*   You started serverless, but are now running into [AWS Lambda Limits](http://docs.aws.amazon.com/lambda/latest/dg/limits.html).
*   You have an existing application server and would like a gradual migration path to move routes to serverless.
*   You’d like to try out an APIGateway to Lambda front end, and an SQS-backed worker pool in the background.
*   You’d like to use the new [Application Load Balancer](https://aws.amazon.com/blogs/aws/new-aws-application-load-balancer/) and HTTP/2 or WebSockets together with Lambda.
*   You’d like to run #serverless in Docker in #serverless and start a Twitter war.

Beginning with [Sparta 0.8.0](https://github.com/mweagle/Sparta), you can now:

1.  Define your service’s responsibility in a single Go application.
2.  Deploy, in a single CloudFormation operation:

*   Specific functions to AWS Lambda.
*   A Docker-ized version of the application to ECS.

Sparta can build a statically-linked version of your Go application, create a Docker image from it, push it to the [EC2 Container Registry](https://aws.amazon.com/ecr/), and provision a CloudFormation stack with the new image.
![image](/posts/migrated/2016-08-25_make-sharpserverless-servers-great-again/images/1.png#layoutTextWidth)
Serverless is a powerful new model for rapid and scalable service development. Sparta v0.8.0 enables you to make the transition at your own pace and with confidence that you can seamlessly switch topologies as conditions evolve.

See [SpartaDocker](https://github.com/mweagle/SpartaDocker) for an example as well as the Sparta [documentation](http://gosparta.io). In the next post I’ll walk through how Sparta supports deploying to Docker/ECS and serverless from the same codebase, at the same time.### Sparta 0.8.0 — Make #serverless Servers Great Again

![image](/posts/migrated/2016-08-25_make-sharpserverless-servers-great-again/images/2.png#layoutTextWidth)
