---
title: "Serverless, Serverfull, and Weaving Pipelines"
author: "Matt Weagle"
date: 2017-09-28T02:30:16.110Z
lastmod: 2023-02-18T21:57:34-08:00
tags: ["serverless", "sparta"]
description: ""
categories: 
    - Medium
subtitle: "Serverless computing is a hot topic these days. Depending on the tweet-weather, it’s a cloud computing revolution, a ripe source of…"

image: "/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/4.jpeg"
images:
 - "/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/1.jpeg"
 - "/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/2.jpeg"
 - "/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/3.jpeg"
 - "/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/4.jpeg"
 - "/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/5.png"
---

![image](/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/1.jpeg#layoutTextWidth)


Serverless computing is a hot topic these days. Depending on the tweet-weather, it’s a cloud computing revolution, a ripe source of lock-in, a net cost savings or an unpredictable cost driver, a stepping stone to event-sourced applications, the commoditization of containers, and a host of other things. Or maybe it’s just a full circle return to the glory days of Perl and [cgi-bin](https://medium.com/adobe-io/2017-will-be-the-year-of-the-cgi-bin-err-serverless-f5d99671bc99).

Putting aside the tweet-bait, there’s a general agreement that running a “serverless service” definitely doesn’t signal the End of Ops. Operations is characterized by a [never-ending evolution of responsibilities, values, and organizational accountabilities](https://charity.wtf/2016/05/31/wtf-is-operations-serverless/) (you should listen to [Charity Majors](https://medium.com/u/5587d135a397)), rather than a specific tool or technology. While serverless might eliminate the need for OS-configuration management tools, it doesn’t eliminate the need for configuration management in general…sadly.

Chief among operational concerns is the ability to safely, quickly, and transparently deliver and sustain production functionality. Minimizing cycle time — the time it takes from code to get from development to customer availability — is one of the hallmarks of Lean. This typically manifests itself as a build pipeline of varying degrees of sophistication ([Spotify](https://continuousdelivery.com/2014/02/visualizations-of-continuous-delivery/), [Lean Enterprise](https://www.amazon.com/Lean-Enterprise-Performance-Organizations-Innovate/dp/1449368425)). Adrian Cockcroft puts continuous delivery at the [center](https://dzone.com/articles/key-takeaways-adrian-cockrofts) of an organization’s ability to react. It’s the OODA engine. If you can’t reliably deploy, then microservices might not be the first thing you want to try.

![image](/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/2.jpeg#layoutTextWidth)


#### Serverfull Serverless

Serverless services make the need for a structured pipeline even more acute, precisely _because_ the deployment friction is so low. After the thrill of sub-three minute deploys from a laptop has worn off, my PagerDuty-addled mind thinks of the night terrors. I have visions of bleary-eyed triage sessions, trying to diagnose a failure, where there is literally no host to SSH into and un-versioned functions only exist for 42ms. I was already despondent thinking about configuration, and now this.

Serverless-based services benefit from CI/CD pipelines as much as non-serverless ones do. Despite these advantages, I’m not really excited about spinning up a server-clinging build cluster to provision an ephemeral service. What I can do instead is leverage related cloud-based services that are effectively serverless from my perspective. Specifically, I can make my serverless service become more [servicefull](https://www.slideshare.net/ServerlessConf/patrick-debois-from-serverless-to-servicefull). And with [Sparta 0.20.0](https://github.com/mweagle/Sparta/blob/master/CHANGES.md#v0200), it’s possible to create a serverless service that provisions its own CI/CD pipeline using a combination of:

*   CodePipeline
*   CodeBuild
*   CloudFormation

Sparta 0.20.0 adds CodePipeline [CloudFormation](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-cfn-artifacts.html) artifact generation to support parameterized stack definitions and promotion.

Serverless! Hoist thyself to the cloud! (But please do it in a standardized way. This isn’t the place for improvisation.)

#### Details

You can add CI/CD provisioning to your Sparta service using these steps:

1.  Define the CodePipeline environments
2.  Register a custom application command
3.  Define the CodePipeline stack with CloudFormation
4.  Add a [_buildspec.yml_](http://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-syntax) file to your repo
5.  Execute the custom command

**Define the CodePipeline Environments**

The first step is to define the set of environments that represent CodePipeline [stages](http://docs.aws.amazon.com/codepipeline/latest/userguide/tutorials-four-stage-pipeline.html). Each stage produces an independent CloudFormation stack that is created via [template configuration files](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-cfn-artifacts.html#w2ab2c13c15c15). The [SpartaCodePipeline](https://github.com/mweagle/SpartaCodePipeline) sample defines a two-stage pipeline: _test_ and _production._




The _RegisterCodePipelineEnvironment_ call accepts an environment name and a set of environment variables that are made available to stacks running in that stage.

Your [lambda function](https://github.com/mweagle/SpartaCodePipeline/blob/master/main.go#L27) references the environment variable in the normal way.




**Register a Custom Command**

Next you register a [cobra.Command](https://github.com/spf13/cobra) with Sparta so that your service can intercept the standard command-line handler. This _provisionPipeline_ command also needs some additional command-line values:

*   _pipeline_: The CodePipeline name
*   _repo_: The GitHub URL hosting your code (https://github.com/mweagle/SpartaCodePipeline)
*   _oauth_: A GitHub Auth token to supply to CodePipeline to support GitHub notifications
*   _s3Bucket_: An S3 bucket name where the CodePipeline CloudFormation template should be uploaded. CodePipeline intermediate artifacts use an S3 bucket defined by the template.
![image](/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/3.jpeg#layoutTextWidth)


**Define the CodePipeline…Pipeline**

Then you define a CloudFormation stack that includes the two-stage CodePipeline resource. Our sample application pushes this into [pipeline.go](https://github.com/mweagle/SpartaCodePipeline/blob/master/pipeline/pipeline.go) to separate it from the normal AWS Lambda execution code.

Even discounting the Go overhead, there is a tremendous amount of configuration needed. The sample app also uses overly permissive IAM roles to help minimize the code size. Despite this, it’s still a lot of configuration. Making this simpler to express is something I think could be significantly improved and would encourage more sophisticated pipelines.

Once everything is lined up, you can use existing Sparta functionality (_ConvergeStackState_) to handle marshaling the template, uploading, and appropriately provisioning or updating the existing pipeline.

**Add buildspec.yml**

Before you can provision the pipeline and actually trigger the build, add a _buildspec.yml_ file to tell [CodeBuild](https://aws.amazon.com/codebuild/) how to handle a Sparta application.

Your [_buildspec.yml_](https://github.com/mweagle/SpartaCodePipeline/blob/master/buildspec.yml) __ file’s lifecycle hooks can be broken down into the following lifecycle stages:

*   The `pre_build` step fetches the application dependencies, moves them to the appropriate location, and installs some OS utilities.
*   The `build` step command is similar to a normal provision command, but includes the **codePipelinePackage** flag. This instructs Sparta to produce a ZIP file that supports [CodePipeline &amp; CloudFormation](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-basic-walkthrough.html). You can inspect the ZIP contents by running the _provisionPipeline_ command locally with the`--noop` flag.
*   The `post_build`step unzips the default archive so that the CodePipeline artifacts aren’t double-ZIP’ped. See the [docs](http://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-syntax) for more info.

**Execute the Command**

Now that you have configured everything, it’s time to run the new _provisionPipeline_ command and wait for the inevitable successful completion.

It’ll work.

With the pipeline provisioned and subscribing to GitHub notifications, every push to the master branch will trigger a new pipeline execution. No servers involved.

The sample two-stage pipeline includes two manual approval steps to mimic a QA pass. Assuming both steps are approved, you will see two new Sparta stacks in your account, one for each envrionment:

*   Test-SpartaCodePipeline-master
*   Prod-SpartaCodePipeline-master

You can extend this pipeline with [Rollback Triggers](https://aws.amazon.com/about-aws/whats-new/2017/08/aws-cloudformation-adds-rollback-triggers-feature/) to automate the manual approval, automatically include a Sparta-provided [CloudWatch dashboard](https://github.com/mweagle/Sparta/blob/master/CHANGES.md#v0130) for each stack, or even support [builds for PRs](https://aws.amazon.com/about-aws/whats-new/2017/09/aws-codebuild-now-supports-building-github-pull-requests/).

#### Two Fewer Things

One of the drivers behind adding CI/CD support is something I only briefly mentioned. Take a look at the Sparta lambda function again:




This release migrates to Go’s standard`[http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)` as AWS Lambda targets! Using the standard signature allows you to chain [middleware](https://github.com/justinas/alice) for your AWS Lambda functions. Formal arguments that used to be in the [sparta.LambdaFunc](https://godoc.org/github.com/mweagle/Sparta#LambdaFunction) are now available in the request context:




The previous [sparta.LambdaFunction](https://godoc.org/github.com/mweagle/Sparta#LambdaFunction) type is officially deprecated and will be removed in a future release. Legacy function signatures remain supported and are noted by a WARN log message:
`WARN[0009] DEPRECATED: sparta.LambdaFunc() signature provided. Please migrate to http.HandlerFunc() Name=main_transformImage`
![image](/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/4.jpeg#layoutTextWidth)


### Summary

Beyond enabling CI/CD support, Sparta 0.20.0 is a significant release for other reasons. The cold-start bootstrapping penalty is reduced and all NodeJS/Python/Go calls now use [protocol buffers](https://github.com/mweagle/Sparta/blob/mweagle/0.20.0/proxy/proxy.proto) for performance and extensibility. See the [change notes](https://github.com/mweagle/Sparta/blob/master/CHANGES.md) for the complete rundown.

At a higher level, this release continues the effort to bring functional and operational needs to a common place. While I was working on this feature, I listened to the [Software Defined Talk episode](http://www.softwaredefinedtalk.com/106) about monitoring. That episode included the epic statement “Monitoring sucks, no you suck.” I heard this while riding the train and laughed out loud. Afterwards, I was reminded of [Kelsey Hightower](https://medium.com/u/9e783a6f12f6)’s [excellent _healthz_ talk](https://vimeo.com/173610242) and the resiliency patterns [Michael Nygard](https://medium.com/u/b1de9e016291) discusses in his [_Release It_](https://www.amazon.com/Release-Production-Ready-Software-Pragmatic-Programmers/dp/0978739213)_!_.

In each case, the best solution is to integrate the non-functional concerns into the codebase itself, rather than treating the code as an inviolate black box. In production, the lines between what is functional and operational are blurry at best, and having [environmental sympathy](https://read.acloud.guru/environmental-sympathy-e6e2f4933b90) for where and how code will execute helps mitigate downstream integration and triage costs.

Given how much business functionality can be expressed in serverless and how intrinsically intertwined that is with cloud infrastructure, I think there is a great opportunity to reduce the barriers between operations and development. You can then shift your focus away from perceived borders and move towards weaving different services together in order to deliver value more reliably, transparently, and collaboratively.

![image](/posts/migrated/2017-09-28_serverless-serverfull-and-weaving-pipelines/images/5.png#layoutTextWidth)


Photo by [Igor Ovsyannykov](https://unsplash.com/photos/w_nh1ECO7QY?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText) on [Unsplash](https://unsplash.com/?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText)
