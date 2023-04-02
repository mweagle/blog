---
title: "Sparta v1.5.0— The Observability Edition"
author: "Matt Weagle"
date: 2018-10-23T05:28:43.385Z
lastmod: 2023-02-18T21:57:39-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "This month marks the three (!) year anniversary of my work on Sparta. Sparta is a framework that transforms a go binary into a GitOps…"

image: "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/1.jpeg"
images:
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/1.jpeg"
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/2.jpeg"
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/3.jpeg"
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/4.jpeg"
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/5.jpeg"
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/6.gif"
 - "/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/7.png"


aliases:
    - "/sparta-v1-5-0-the-observability-edition-a8257eeb11d6"

---

![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/1.jpeg#layoutTextWidth)


This month marks the three (!) year anniversary of my work on [Sparta](https://gosparta.io). Sparta is a framework that transforms a go binary into a GitOps friendly, self-deploying, and operationally aware application that targets AWS Lambda as its execution environment.

Given a function similar to this and a bit of bootstrapping code:
`func helloWorld() (string, error) {
  return &#34;Hello World&#34;, nil
}`

With a single [mage](https://magefile.org/) command you can compile, package, upload, and create a CloudFormation managed microservice. (The log is more colorful in person).




See the [SpartaHelloWorld](https://github.com/mweagle/SpartaHelloWorld) project for more details.

One of the guiding principles for Sparta’s development is that that functional and non-functional requirements should be uniformly expressed. [Metaparticle.io](https://metaparticle.io/about/) does something similar at a library level. Languages and the assumptions they tacitly embed often foster distinct language communities that can come to see one another as “different”. This perceived distance can hinder collaboration, which from a customer perspective, often means that things “don’t work”. Customers don’t care that a service self-reports it’s working at 5 9s, but the bash to Chef migration caused the 12-factor config to drift a bit and now a [critical environment variable is missing](https://blog.acolyer.org/2016/11/29/early-detection-of-configuration-errors-to-reduce-failure-damage/) and in fact, as far the customer is concerned the service is not working at all. That last part is something customers most definitely care about — until they don’t.

On a new project it might start out being accurate that infrastructure is relatively rigid and the business logic seems endlessly fluid. (Those customers, what exactly **do** they want?) **** However, perhaps a few years later the service has taken off and the business logic is virtually in lock-down, the team disbanded. Only now the organization is moving to new regions, providers, or execution environments. Or maybe even eliminating a piece of infrastructure entirely in favor of a [Service Full](https://www.slideshare.net/jedi4ever/from-serverless-to-service-full-how-the-role-of-devops-is-evolving) migration. What was fluid is no longer so and what was presumed fixed is now in a constant state of churn. Rates of change…change.

Promoting operational responsibilities to the same tier as business logic opens up the possibility of making more adaptive and internally consistent service deployments. There’s a possibility that a deployment can close over all its typically dangling dependencies and ensure that the ids, metrics, alerts, functions, secrets, dashboards, per-environment behaviors and everything else over in JSON/YAML/XML/ConfigLand can be expressed in a uniform way.

Which sets the stage for Sparta 1.5.0 — The Observability Edition. At the risk of enraging both _serverless_ and _observability_ purists, I’m going to stick with the general idea of “How can I gain a better understanding of my existing service.” There are others who can speak much better to observability itself including [JBD](https://medium.com/u/1737b4e67578), [Charity Majors](https://medium.com/u/5587d135a397), and [Cindy Sridharan](https://medium.com/u/87c8c84f24b1). I’m limiting things to the case of how a Sparta service can provide more transparency using the same programming constructs used to define the service behavior.

![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/2.jpeg#layoutTextWidth)


#### Metrics

To round out your service’s core business logic, Sparta provides an opportunity to decorate the CloudFormation template that defines your infrastructure. For example, you can provision a CloudWatch Dashboard using a [DashboardDecorator](https://godoc.org/github.com/mweagle/Sparta/decorator#DashboardDecorator) to produce an application-centric view of your service (similar to the recently announced [Applications](https://aws.amazon.com/about-aws/whats-new/2018/08/aws-lambda-console-enables-managing-and-monitoring/) view).

The past few Sparta releases have extended those capabilities to include support for:

*   [CloudWatchErrorAlarmDecorator](https://godoc.org/github.com/mweagle/Sparta/decorator#CloudWatchErrorAlarmDecorator): create and associate a CloudWatch Alarm that’s triggered by a configurable number of AWS Lambda errors over a given period.
*   [RegisterLambdaUtilizationMetricPublisher](https://godoc.org/github.com/mweagle/Sparta/aws/cloudwatch#RegisterLambdaUtilizationMetricPublisher): register a periodic task that publishes container level metrics to CloudWatch. Define custom [metric dimensions](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/cloudwatch_concepts.html#Dimension) including your deployment’s `BuildID` (compile time) or `InstanceID` (runtime) values, or both. Among the metrics included are _CPUPercent_ and _DiskUsedPercent…_not that I’ve had issues with hitting the [512 MB limit](https://docs.aws.amazon.com/lambda/latest/dg/limits.html) or anything.
![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/3.jpeg#layoutTextWidth)


#### Logs

Logging is a critical observability capability for any microservice deployment. Having a durable and searchable centralized logging store is especially important for serverless-based solutions. [Yan Cui](https://medium.com/u/d00f1e6b06a2)’s [Centralized logging for AWS Lambda](https://hackernoon.com/centralised-logging-for-aws-lambda-b765b7ca9152) is an excellent introduction to the topic.

Sparta has always supported structured logging via [logrus](https://github.com/sirupsen/logrus) and aggregation to CloudWatch Logs. However as Yan discusses, there are other options that may offer more features or work better within your organization. To support this use case, Sparta now includes a [LogAggregatorDecorator](https://godoc.org/github.com/mweagle/Sparta/decorator#LogAggregatorDecorator) that forwards all referenced log statements to a Kinesis Stream for asynchronous processing.

Ensuring Kinesis has has a copy of all log statements doesn’t really help on the observability front though. Even better than having the logs is the ability to effectively search them. The [SpartaPProf](https://github.com/mweagle/SpartaPProf/blob/master/pprof/googlePProf.go#L102) example turns the Kinesis stream into StackDriver log events and forwards them to Google Stackdriver. Use [Stackdriver logging](https://cloud.google.com/logging/) to search your AWS Lambda logs!

![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/4.jpeg#layoutTextWidth)
Searching AWS Lambda logs in Google Stackdriver #multicloud



#### Profiling

While profiling doesn’t typically constitute one of the three [pillars of observability](https://cengizhan.com/3-pillars-of-observability-8e6cb5434206), it’s still often helpful to understand your service’s performance. Leveraging go’s ability to [expose profiling data](https://golang.org/pkg/runtime/pprof/), Sparta 1.5.0 provides facilities to make it straightforward to send AWS Lambda extracted profiling information to [Google Stackdriver Profiling](https://cloud.google.com/profiler/) for visualization. Using Google credentials stored in [AWS Systems Manager Parameter Store](https://github.com/mweagle/ssm-cache), you can now install a profiling task to really understand why your function has crossed the 100ms billing increment level.

See the [SpartaPProf](https://github.com/mweagle/SpartaPProf/blob/master/pprof/googlePProf.go#L152) sample project for a complete example.

![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/5.jpeg#layoutTextWidth)
Stackdriver profiling of AWS Lambda executions #multicloud



#### Additional Treats

There are a host of other improvements and bug fixes, including:

*   Creating your own `go` CloudFormation [CustomResources](https://gosparta.io/reference/custom_lambda_resources/)! Want to do something that isn’t supported in CloudFormation? Want to call out to a third party API as part of your service’s normal lifecyle? Custom resources provide an “escape hatch” for those times when nothing but custom code will get the job done.
*   Publishing an S3 artifact as part of your service’s lifecycle with a [S3ArtifactPublisherDecorator](https://godoc.org/github.com/mweagle/Sparta/decorator#S3ArtifactPublisherDecorator). This is a great fit for those times when you want your service to leave some sort of metadata receipt in a bucket.
*   New [archetype](https://godoc.org/github.com/mweagle/Sparta/archetype) constructors to eliminate much of the boilerplate around typical serverless patterns such as subscribing to S3 or DynamoDB events.
*   Putting a [CloudFront Distribution](https://godoc.org/github.com/mweagle/Sparta/decorator#CloudFrontSiteDistributionDecorator) and custom domain in front of your S3-backed static site
*   Defining Validator [WorkflowHooks](https://godoc.org/github.com/mweagle/Sparta#WorkflowHooks) that receive an immutable version of the complete CloudFormation template. Validator hooks are useful to define team or organization-wide stack policies (eg: prevent `Resource: *` in IAM policies when possible). Those policies can be distributed as standard go packages.
*   Using [magefile tasks](https://godoc.org/github.com/mweagle/Sparta/magefile) for cross-platform friendly actions such as provisioning a service or applying a tool to all `*.go` files in your source tree. Moving from Makefiles to mage makes it much easier to support mixed *nix/Windows teams. You can largely copy the new _magefile.go_ files across projects and they will Just Work. Here’s an [example](https://github.com/mweagle/SpartaHelloWorld/blob/master/magefile.go).

#### Get On Board the WASM Train

![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/6.gif#layoutTextWidth)


Finally, go 1.11 added experimental support for [WebAssembly](https://github.com/golang/go/wiki/WebAssembly). This compiler target opens up the possibility of using go to write front end code as well!

As the Sparta provisioning lifecycle supports running `go:generate` as part of the cross compilation step, it’s now feasible to compile go to WASM and deploy that to your static S3 site as part of a single provision step. See the [SpartaWASM](https://github.com/mweagle/SpartaWASM) repo for an example.

And full disclosure, the WASM train is desperately in need of more rail. Pull Requests are most definitely appreciated!

### Wrapping Up

The Sparta 1.5.0 release is designed to streamline the developer experience and provide facilities to better understand the state and shortcomings of your Sparta service. While the speed of serverless development and deployment is addictive, I think it’s also important to ensure that the longer-term operational aspects of your service enjoy the same level of integration and expressiveness as your core logic.

And maybe even more importantly, the WASM train is coming and I look forward to seeing many of you at the station with me!

Comments appreciated and if you run into any issues, please open an issue at the [Sparta repo](https://github.com/mweagle/Sparta). Additional documentation is available at [https://gosparta.io](https://gosparta.io.) and full change notes at [CHANGES](https://github.com/mweagle/Sparta/blob/master/CHANGES.md). Build something awesome.
![image](/posts/2018-10-23_sparta-v1.5.0-the-observability-edition/images/7.png#layoutTextWidth)


Credits

*   [Kyler Boone](https://unsplash.com/photos/4bY-VkQ0UZQ)
*   [Stephen Dawson](https://unsplash.com/photos/qwtCeJ5cLYs)
