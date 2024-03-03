---
title: "Build an S3 website with API Gateway and AWS Lambda for Go using Sparta"
author: "Matt Weagle"
date: 2018-01-22T04:29:17.440Z
lastmod: 2023-02-18T21:57:36-08:00
tags: ["serverless", "sparta"]
description: ""
categories: 
    - Medium
subtitle: "Deploy a complete microservice in less than 100 lines of code"

image: "/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/1.png"
images:
 - "/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/1.png"
 - "/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/2.jpg"
 - "/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/3.jpg"
---

#### Deploy a complete microservice in less than 100 lines of code

![image](/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/1.png#layoutTextWidth)


While AWS Lambda has always supported the ability to run [arbitrary binaries](https://aws.amazon.com/blogs/compute/running-executables-in-aws-lambda), enabling a response to Lambda requests had always proven difficult — requiring some type of shim layer to proxy requests through one of the supported runtimes.

Over the past few years, I’ve been working on a Go framework that produces self-deploying AWS Lambda-powered microservices. [Sparta](http://gosparta.io) offered both NodeJS (via a Go binary) and Python (via [cgo](https://golang.org/cmd/cgo) shims) — but each approach suffered from both runtime performance penalties and developer complexity.

The recent announcement that AWS Lambda officially supports Go eliminates these limitations, and makes Go an excellent option for writing AWS Lambda services.

> [](https://twitter.com/timallenwagner/status/953052636592316418)


In this post, I’ll outline the developer changes for AWS and Sparta — and end with an overview of how to deploy a complete service that includes a static HTML site using Amazon S3, an API Gateway CORS-enabled HTTP resource, and an AWS Lambda Go function.

In less than 100 lines of custom Go code, Sparta can help you deploy a complete microservice composed of multiple AWS services — managed exclusively through CloudFormation.

#### **Requirements**

To build and deploy the service, you’ll need to:

1.  Install [go](https://golang.org/doc/install).
2.  `go get -u -v github.com/mweagle/SpartaHTML` to fetch the package.
3.  Navigate to the _SpartaHTML_ directory and fetch the dependencies via `go get -u -v ./...`
4.  Configure your AWS credentials per the [FAQ](http://gosparta.io/docs/faq). Note that the target region is set by the AWS_REGION environment variable.
5.  Ensure you have access to an S3 bucket in the region you plan to deploy the service.

#### **An overview of Sparta**

Before discussing the Sparta programming model changes, let’s first do a quick walkthrough of a _Hello World_ Sparta-based application.

Sparta is a framework that accepts a set of user-defined functions and creates a self-deploying binary that:

*   [Cross-compiles](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) your Go source code into an AWS Lambda compatible executable.
*   Creates a CloudFormation template for deployment.
*   Zips and uploads all assets to S3.
*   Uses CloudFormation APIs to perform the requested operation (`provision`, `delete`, etc.).

Let’s take a look at a pre-1.0 Sparta _Hello World_ application that defines a single function, `helloWorld`, that is supplied to Sparta via the `lambdaFunctions` slice.




Sparta transforms the set of user functions into a self-deploying binary that exposes a command-line interface. To see which options are available, just type `go run main.go`:




To provision a CloudFormation Stack for this service, use the command:

`go run main.go provision --s3Bucket my-S3-bucketName`

In the parameters, `my-S3-bucketName` is an S3 bucket that you can access with your credentials. By default, all Sparta operations are marshaled to a CloudFormation template — including configuring all supporting AWS services and any [custom resources](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources-lambda.html) necessary to complement the core functionality.

#### **Out with the old**

What’s worth pointing out in the previous section is the somewhat conspicuous Go HTTP-style function signature:

`func(w http.ResponseWriter, r *http.Request)`.

Over the course of development, Sparta has supported both this signature and the following one as valid AWS Lambda targets:
`type LambdaFunction func(*json.RawMessage,
    *sparta.LambdaContext,
    http.ResponseWriter,
    *logrus.Logger)`

Both signatures were concessions to the need for a proxying tier. When a Lambda request was made, Sparta would marshal back and forth via [protocol buffers](https://developers.google.com/protocol-buffers/) to the Go function exposed by a spawned process.

While this was largely transparent to development, it did not provide the same level of observability or performance as the officially supported runtimes. The newly released Sparta 1.0 now supports the official AWS Lambda Go programming model.

**The AWS Lambda Go programming model
**The [AWS Lambda Go SDK](https://github.com/aws/aws-lambda-go) supports the following function signatures as valid Lambda handlers:

*   `func()`
*   `func() error`
*   `func(TIn) error`
*   `func() (TOut, error)`
*   `func(context.Context) error`
*   `func(context.Context, TIn) error`
*   `func(context.Context) (TOut, error)`
*   `func(context.Context, TIn) (TOut, error)`

The `TIn` and `TOut` parameters represent [encoding/json](https://golang.org/pkg/encoding/json) un/marshallable types. The SDK uses reflection at startup to determine which function signature to invoke for the duration of the Lambda handler lifespan. A Go binary, by default, can only register a single handler.

> [](https://twitter.com/ashleymcnamara/status/914992633868640257)


#### **In with the new**

To maintain full compatibility with the AWS Lambda Go runtime, Sparta now accepts only AWS Lambda Go SDK signature-compatible functions.

In the prior example, only two changes are required.

First, the function definition itself:




Secondly, the `sparta.HandleAWSLambda` call was updated by removing the `http.HandlerFunc` type:




These are the only changes required to upgrade this application to Sparta 1.0 and the officially supported [Lambda Go runtime](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model.html). During compilation, Sparta will validate all provided Lambda function signatures.

If a caller inadvertently supplies an unsupported signature:




Sparta will report an error during `provision`:
`ERRO[0000] Lambda function (Hello World) has invalid returns: handler returns a single value, but it does not implement error
exit status 1`

#### **Deploying a website**

Now that we’ve taken a look at the changes in Sparta 1.0, let’s walk through a complete service. The full source code is available at [SpartaHTML](https://github.com/mweagle/SpartaHTML).

**#1 Define the Lambda function
**The first step is to define a Go function that adds a property to the incoming request and sends the consolidated response back to the requestor:




**#2 Include a website
**The second step is to include a [static S3 site](https://docs.aws.amazon.com/AmazonS3/latest/dev/WebsiteHosting.html) deployment as part of the microservice.
`s3Site, s3SiteErr := sparta.NewS3Site(&#34;./resources&#34;)`

The code above instructs Sparta to zip everything in the _./resources_ directory and deploy it to a dynamically created S3 bucket as part of the service creation. The S3 bucket permissions will also be modified to allow public access.

**#3 Define the API Gateway
**To make the Go function publicly available, the service must provision an API Gateway REST API. This is done in the following snippet where we also limit CORS access to the dynamically provisioned S3 bucket’s [WebsiteURL](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket.html) value:




The Lambda definition is then linked to the API Gateway _/hello_ URL path:




**#4 Delegate to Sparta
**With the major pieces in place, the final step is to supply the Lambda function definitions and associated AWS service configurations to Sparta’s primary entrypoint:




**#5 Provision the Service
**The service is now fully defined and can be provisioned to a supported [AWS region](https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services) with the following command:
`go run main.go provision --s3Bucket my-S3-bucketName`

As a result of the `provision` step, you should see output similar to the following:




**#6 Visit the Website!
**The most important section in the log output is the _Stack Outputs_ section that includes both the S3 website and API Gateway URLs. If you visit the **S3 Website URL** value, you should see a congratulations page:

![image](/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/2.jpg#layoutTextWidth)


Click the **Get Started** button to scroll the page, where you can then click the **Try It!** button to call your Go Lambda function and see the results:

![image](/posts/migrated/2018-01-22_build-an-s3-website-with-api-gateway-and-aws-lambda-for-go-using-sparta/images/3.jpg#layoutTextWidth)


In approximately 90 seconds we have deployed a completely self-contained CloudFormation-defined service — composed of a Lambda function, an API Gateway instance, and an S3-backed website in under 100 lines of custom Go code. _Congratulations!_

#### **Some final words on Sparta**

AWS Lambda is a transformational service that positions you to deliver more resilient and scalable microservices while also dramatically reducing development cycle time. It is often used as a critical component in larger [event-based architectures](https://softwareengineeringdaily.com/2017/11/14/serverless-event-driven-architecture-with-danilo-poccia) that are composed of multiple AWS services.

Sparta is designed to make it easier to programmatically assemble those individual services into higher-order microservices.

Sparta offers an extensive feature set, including registering for [S3 events](http://gosparta.io/docs/eventsources/s3), locally viewing [Go runtime profile data](https://golang.org/pkg/runtime/pprof/) across your Lambda containers, and even creating [AWS Step Functions](https://github.com/mweagle/SpartaStep).

These features leverage the rich ecosystem of Go libraries and tools to produce a single service binary. This executable is deployed and managed through CloudFormation and makes a very performant AWS Lambda handler due to its relatively small size and low runtime overhead.

Ready to make “AWS Lambda Go”? Visit the [AWS Go announcement post](https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/), the AWS Lambda Go SDK [documentation](https://godoc.org/github.com/aws/aws-lambda-go) or the Sparta [documentation](http://gosparta.io) to learn more and start building today!
