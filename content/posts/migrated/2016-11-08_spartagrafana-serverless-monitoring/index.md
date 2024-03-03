---
title: "SpartaGrafana — Serverless Monitoring"
author: "Matt Weagle"
date: 2016-11-08T15:50:49.504Z
lastmod: 2023-02-18T21:57:29-08:00
tags: ["serverless", "sparta"]
description: ""
categories: 
    - Medium
subtitle: "One of the most contentious and pedantically debated points around serverless is the term itself."

image: "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/1.png"
images:
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/1.png"
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/2.gif"
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/3.png"
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/4.jpeg"
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/5.png"
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/6.png"
 - "/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/7.gif"

---

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/1.png#layoutTextWidth)


One of the most contentious and pedantically debated points around serverless is the term itself.

Luckily, I’m not going to talk about that. Let’s just call it [Jeff](https://serverless.zone/serverless-is-just-a-name-we-could-have-called-it-jeff-1958dd4c63d7#.w7x9joip9).

What seems less controversial is that serverless represents an opportunity to quickly develop truly cloud-native applications. Applications that begin to intrinsically incorporate broader operational aspects into their core codebase. If serverless is about [super-advanced cloud powers](https://medium.com/@PaulDJohnston/is-serverless-just-another-way-of-saying-super-advanced-cloud-1f652357620f#.s5edlbatc), how might the avengers of features, operations, security, capacity, and other apparently disjoint domains come together in a single [repository](https://serverless.zone/serverless-deploy-thyself-6085412f5393), with a bias towards a shared language and a common goal?

### Graph All The Things

At the first [Seattle Serverless Meetup](http://www.meetup.com/Serverless-Seattle/), Rob Gruhl from Nordstrom presented a preview of their [serverless-artillery](https://github.com/Nordstrom/serverless-artillery) load-testing tool. He gave an excellent talk and one of the things I took away from his presentation was the incredible level of observability their team had attained. The ability to visualize the increase in load against their system, with minimal latency, was very powerful. Rob and the team had done this by instrumenting their code and publishing the results to [Grafana](http://grafana.org/).

Not to mention it looked really cool.

While watching the presentation, I made a mental note to somehow add support for spinning up a Grafana instance as part of a [Sparta](http://gosparta.io) application. This of course meant extending Sparta beyond serverless, into that vast and well-traveled territory of…servers.

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/2.gif#layoutTextWidth)


### Servers?

Serverless is a good fit for many workloads and stages of product development, but running Grafana (a SaaS version is [coming](http://grafana.org/hosting/)!) isn’t one of them. However, the recent addition of CloudFormation [cross-stack references](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/walkthrough-crossstackref.html) made it conceivable to both decouple and bridge the two worlds.

The [SpartaGrafana](https://github.com/mweagle/SpartaGrafana) application demonstrates how to use the new Sparta [WorkflowHooks](https://godoc.org/github.com/mweagle/Sparta#WorkflowHooks) feature to provision a completely independent CloudFormation stack that supports a single EC2 instance running Grafana to which AWS lambda instances publish metrics.

#### Grafana Stack Definition

The Grafana stack definition is fully specified by [_grafana/grafana.go_](https://github.com/mweagle/SpartaGrafana/blob/master/grafana/grafana.go) __ which:

*   Defines the CloudFormation stack, including a CloudFormation export that publishes the EC2 _PublicDnsName_.
*   Includes a CloudInit shell script that downloads Grafana, InfluxDB, and scripts to bootstrap the InfluxDB datasource and Grafana _Sparta Hello World_ dashboard.

The stack status is ensured by [ConvergeStackState](https://godoc.org/github.com/mweagle/Sparta/aws/cloudformation#ConvergeStackState), which does most of the heavy AWS lifting (marshaling, uploading to s3, waiting for the stack to complete, etc.).

#### Behind the Serverless Curtain

Once the Grafana stack is created/updated in `PostBuildHook`, the normal Sparta [provisioning workflow](http://gosparta.io/docs/intro_example_details/#provisioning-flow) continues. This deploys a HelloWorld lambda function together with a single-resource API Gateway stage.




The most interesting aspect of the HelloWorld lambda function is that it integrates with the [go-metrics](https://github.com/rcrowley/go-metrics) library. Each lambda invocation increases a local counter via `helloWorldCounterMetric.Inc(1)`. The lambda instance is able to discover the Grafana EC2 DNS name by first `Fn::Import`ing the value:




which is then looked up at initialization time to create the [InfluxDB publisher](https://github.com/vrischmann/go-metrics-influxdb):




Using this information, each lambda container instance sets up a publishing loop that includes a randomly generated tag value to identify the lambda container instance.

The benefit of the `Fn::Import` statement is that it prevents “fat-finger” deletions via the AWS Console. It’s not possible to delete the _Grafana_ stack while the _SpartaGrafanaPublisher_ stack is active.

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/3.png#layoutTextWidth)


#### Results

Provisioning _SpartaGrafana_ creates two CloudFormation stacks. The CloudFormation Outputs for the _Grafana_ stack include the URL to log in to the Grafana console (`admin/admin` is the default _username_ and _password_).

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/4.jpeg#layoutTextWidth)


After logging in to Grafana, navigate to the pre-created _Sparta Hello World_ dashboard. This dashboard tracks a single metric (`helloWorldCounterMetric`) that graphs how frequently the lambda function was invoked.

We’ll use [https://goad.io](https://goad.io) to generate some load against our lambda function. First we need to find the API-Gateway URL for our lambda function URL via the AWS Console ( API Gateway➤Sparta Grafana➤Stages➤dev➤GET — /hello/grafana). Copy that value and paste it into the Goad web form:

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/5.png#layoutTextWidth)


After a few seconds, you can see AWS Lambda starting to spin up containers to handle the traffic:

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/6.png#layoutTextWidth)


Compared to the previous GIF which only showed 2 containers running, you can see in the image above that AWS Lambda has created 4 containers of varying lifespans to handle this load test.

And thanks to go-metrics, our simple service automatically publishes a host of other metrics that are available for visualization.

![image](/posts/migrated/2016-11-08_spartagrafana-serverless-monitoring/images/7.gif#layoutTextWidth)


### Seamless Serverless

Serverless offers an incredibly easy and potentially transformative way to create highly-scalable and resilient cloud-native applications. However, there are times when the standard (I’m reluctant to say “legacy”) instance-based deployments make sense. Occasionally, a single service may even need to deploy multiple isolated topologies that work together to support a business goal.

And other times they need to work together to produce cool graphs that help you visualize what your serverless application is actually doing.

#### Links

*   Sparta: http://gosparta.io
*   go-metrics: [https://github.com/rcrowley/go-metrics](https://github.com/rcrowley/go-metrics)
*   go-metrics InfluxDB publisher: [https://github.com/vrischmann/go-metrics-influxdb](https://github.com/vrischmann/go-metrics-influxdb)
*   Grafana: [http://grafana.org/](http://grafana.org/)
*   Serverless Slack signup page: [https://wt-serverless-seattle.run.webtask.io/serverless-forum-signup?webtask_no_cache=1](https://wt-serverless-seattle.run.webtask.io/serverless-forum-signup?webtask_no_cache=1)
