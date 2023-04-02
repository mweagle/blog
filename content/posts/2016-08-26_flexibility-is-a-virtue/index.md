---
title: "Flexibility is a Virtue"
author: "Matt Weagle"
date: 2016-08-26T06:45:18.222Z
lastmod: 2023-02-18T21:57:25-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "Sparta v0.8.0 adds a powerful new WorkflowHooks feature that allows service developers the ability to customize the standard provisioning…"

image: "/posts/2016-08-26_flexibility-is-a-virtue/images/1.png"
images:
 - "/posts/2016-08-26_flexibility-is-a-virtue/images/1.png"
 - "/posts/2016-08-26_flexibility-is-a-virtue/images/2.jpeg"
 - "/posts/2016-08-26_flexibility-is-a-virtue/images/3.jpeg"
 - "/posts/2016-08-26_flexibility-is-a-virtue/images/4.png"


aliases:
    - "/flexibility-is-a-virtue-54059d75b1ef"

---

[Sparta](http://gosparta.io/) v0.8.0 adds a powerful new [WorkflowHooks](http://gosparta.io/docs/application/workflow_hooks/) feature that allows service developers the ability to customize the standard provisioning pipeline. This inversion of control makes it possible for a single logical service to be deployed to AWS Lambda **and** the EC2 Container Service…at the same time.

You read that right. **Docker** _and_ **serverless**, working together to create **mass hysteria!**

Or at the very least, working together as part of a single service.

I mentioned some of the reasons you might want to do this in a [previous post](https://medium.com/@mweagle/make-serverless-servers-great-again-260260297d41#.p95alavz6), and in this installment I’m going to walk through how the [SpartaDocker](https://github.com/mweagle/SpartaDocker) sample accomplishes this. It defines a service where Lambda functions push events onto a dynamically created SQS queue, whose messages are consumed by the same application running in ECS.

![image](/posts/2016-08-26_flexibility-is-a-virtue/images/1.png#layoutTextWidth)


In addition to exploring the new Sparta feature I think it’ll become evident how much of the operational burden serverless-style architectures promise to offload.

To the cloud!

#### Create a new application mode

The first step is to create a new application mode that complements the standard Sparta [options](http://gosparta.io/docs/application/commandline/). When the Go application is started in the _sqsWorker_ mode, it will begin consuming messages from an (yet to be defined) SQS provider.




The **sqsListener** function sets up a simple polling loop and writes each received message to CloudWatch Logs.

Since we’re going to be provisioning an AutoScalingGroup that manages 1 EC2 instance, we also need to provide a mechanism to accept the (required) SSH KeyName.




#### Register Workflow Hook

The next step is to define the _PostBuild_ hook that builds the Docker image and pushes it to the EC2 Container Registry.




Sparta provides both [BuildDockerImage](https://godoc.org/github.com/mweagle/Sparta/docker#BuildDockerImage) and [PushDockerImageToECR](https://godoc.org/github.com/mweagle/Sparta/docker#PushDockerImageToECR) to encapsulate the details around building Docker images, authenticating (with retry) to ECR, and pushing your image to the cloud.

As the ECR URL is dynamically determined in part from the _buildID_, it’s stored in the context object so that we can reference it later when building up the CloudFormation template.

#### Docker All the Things

The SpartaDocker Dockerfile is:
`FROM centurylink/ca-certs
ARG SPARTA_DOCKER_BINARY
ADD $SPARTA_DOCKER_BINARY /SpartaDocker
EXPOSE 80
CMD [&#34;/SpartaDocker&#34;, &#34;sqsWorker&#34;]`

The **BuildDockerImage** function compiles and statically links your application’s dependencies (CGO_ENABLED=0). The resulting binary is provided to the Docker build command via the SPARTA_DOCKER_BINARY [Docker build argument](https://docs.docker.com/engine/reference/commandline/build/#/set-build-time-variables-build-arg) and our CMD uses the _sqsWorker_ command line option we defined above.

Building with Docker 1.12 on OSX, the resulting image is 16.66 MB in size.

#### Part 1 — Lambda

The next step is to register a normal Sparta lambda function, together with a [TemplateDecorator](https://godoc.org/github.com/mweagle/Sparta#TemplateDecorator). We’re also going to specify an IAM role so that the Lambda function only has write access to the dynamically provisioned queue (**sqsResourceName**).




It’s the TemplateDecorator that ties everything together.

#### Parts 2 and Beyond

The final step (which is actually multiple steps) is to tie everything together and properly configure and register the ECS cluster. The **helloWorldDecorator** function is by far the largest function in the source. It’s responsible for:

1.  Adding a new **gocf.SQSQueue{}** resource to the CloudFormation template
2.  Provisioning an ECS cluster using the **context[“URL”]** image URL we pushed in the previous step.
3.  Ensuring the ECS cluster can successfully discover the SQS resource

Behind that second step is ~250 lines of configuration files, multiple interdependencies, various AWS service settings, and environment values to properly configure the cluster. I’d encourage you to work through the source if you’re interested in the gory details. I’d also like to thank [Stelligent](https://stelligent.com/2016/05/26/automating-ecs-provisioning-in-cloudformation-part-1/) for their excellent post — the Sparta implementation is largely a port of their writeup.

There are a few points worth calling out that are specific to the SpartaDocker application.

1.  The EC2 IAM Role associated with the ECS host includes privileges to read from the dynamically provisioned SQS instance:



The ECS ContainerTasks inherit this IAM Role which authorizes the Sparta application to access the SQS resource.

2. The dynamically provisioned SQS resource’s properties are provided to the ECS-hosted SpartaDocker application as environment variables:




See the [SQS CloudFormation outputs](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sqs-queues.html) for more details.

### To the Cloud

Once everything is in place, the WorkflowHooks struct is provided to the new [sparta.MainEx](https://godoc.org/github.com/mweagle/Sparta#MainEx) function. Provisioning the service can be done as before:
`make provision`

This will take a bit longer than a serverless-only Sparta deploy, but the end result is a self-contained service that dynamically provisions an SQS resource and supports dynamic discovery so that both a Sparta lambda function and an ECS worker are able to communicate through it.

Using the AWS Console to test the Lambda function:

![image](/posts/2016-08-26_flexibility-is-a-virtue/images/2.jpeg#layoutTextWidth)


You can verify the ECS worker pool processed the message by checking the CloudWatch Logs:

![image](/posts/2016-08-26_flexibility-is-a-virtue/images/3.jpeg#layoutTextWidth)
### Conclusion

While the Serverless model is an ideal fit for many workloads and can significantly reduce both a service’s time to market and ongoing operational costs, it’s not always a perfect match. To support shifting workloads and different topologies, consider adding Docker and ECS deployments to your Sparta application. It does require more configuration and potential monetary costs, but the payoff is a significant increase in flexibility &amp; the ability to exploit new [AWS features](https://aws.amazon.com/blogs/aws/new-aws-application-load-balancer/).

Plus, there’s always the off chance that once Docker and serverless have reconciled, cats and dogs can once again live in peace together.

![image](/posts/2016-08-26_flexibility-is-a-virtue/images/4.png#layoutTextWidth)
