---
title: "Sparta — Behind the Shield"
author: "Matt Weagle"
date: 2016-04-26T14:04:02.793Z
lastmod: 2023-02-18T21:57:21-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "The last post introduced Sparta — A Go Framework for AWS Lambda. In this post I’ll cover some of Sparta’s internals, working off of the…"

image: "/posts/2016-04-26_sparta-behind-the-shield/images/1.png"
images:
 - "/posts/2016-04-26_sparta-behind-the-shield/images/1.png"
 - "/posts/2016-04-26_sparta-behind-the-shield/images/2.gif"
 - "/posts/2016-04-26_sparta-behind-the-shield/images/3.png"


aliases:
    - "/sparta-behind-the-shield-7a6e178f1b72"

---

![image](/posts/2016-04-26_sparta-behind-the-shield/images/1.png#layoutTextWidth)


The last post introduced [Sparta — A Go Framework for AWS Lambda](https://medium.com/@mweagle/a-go-framework-for-aws-lambda-ab14f0c42cb). In this post I’ll cover some of Sparta’s internals, working off of the HelloWorld sample. I’ll touch on the major implementation aspects as of [Sparta 0.5.5](https://github.com/mweagle/Sparta/releases/tag/0.5.5) released on April 22, 2016.

It’s possible that this overview is out of date at the time you’re reading, though. Sparta attempts to provide a consistent and coherent developer-centric view across multiple AWS services and cross-cutting concerns. The implementation reflects what is available on AWS at a given time, and with the rapid rate of AWS change, “code may have shifted during flight”. For the very latest, please see the [official documentation](http://gosparta.io) and [change history](https://github.com/mweagle/Sparta/blob/master/CHANGES.md).

#### Hello World

To get started, let’s review the simple Sparta HelloWorld application:




We’ll dig into what’s going on under the hood with an eye towards understanding how Sparta manages to get your Go function invoked by AWS Lambda.

#### Sparta.Main

Sparta.Main is the primary link between your Go functions and Sparta. It accepts the following arguments (see also the [godoc entry](https://godoc.org/github.com/mweagle/Sparta#Main))

*   **serviceName** (string) — Your application’s name which identifies your application for a given AWS (_AccountId_, _Region_) pair. It is used as the CloudFormation stack name.
*   **serviceDescription** (string)— Optional application description, used in the CloudFormation stack [description field](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-description-structure.html).
*   **lambdaAWSInfos** ([]*sparta.LambdaAWSInfo) — Non-empty slice of Sparta-registered functions that will be transformed into AWS Lambda targets.
*   **api** (*API) — Optional API structure that [provisions an APIGateway](http://gosparta.io/docs/apigateway/) connected to your application.
*   **site** (*S3Site) — Optional S3Site structure that [provisions a dynamically created S3](http://gosparta.io/docs/s3site/) bucket, populates it with static resources, and configures it for [static website hosting](http://docs.aws.amazon.com/AmazonS3/latest/dev/WebsiteHosting.html).

The **serviceName** is particularly significant, as it represents your application’s identity across Sparta operations. When Sparta goes to provision your application, it uses the **serviceName** value to determine if the application has been previously deployed via [DescribeStacks](http://docs.aws.amazon.com/sdk-for-go/api/service/cloudformation/CloudFormation.html#DescribeStacks-instance_method). Using unstable **serviceName** values (eg: _fmt.Sprintf(“MyApp%d”, time.Now().Unix())_ will “trick” Sparta into always cold-provisioning a new CloudFormation stack.

#### *sparta.LambdaAWSInfo

This slice contains the core of your application’s logic: the interrelated set of functionality that defines your service’s (hopefully single) area of responsibility. Each *_LambdaAWSInfo_ element is returned by _sparta.NewLambda_, which accepts

*   **roleDefinition** ([_sparta.IAMRoleDefinition_](https://godoc.org/github.com/mweagle/Sparta#IAMRoleDefinition)) — The IAM role definition under which the Go lambda function will execute. The HelloWorld example defines an empty IAMRoleDefinition value which limits _helloWorld’s_ runtime permissions to a predefined set (see below).
*   **lambdaFunc** ([_sparta.LambdaFunc_](https://godoc.org/github.com/mweagle/Sparta#LambdaFunction)) — The Go function that Sparta invokes in response to an AWS Lambda trigger.
*   **lambdaOptions** ([_sparta.LambdaFunctionOptions_](https://godoc.org/github.com/mweagle/Sparta#LambdaFunctionOptions)) — Optional structure that specifies additional [AWS Lambda execution options](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-function.html).

The _NewLambda_ parameter order underscores the importance of integrating security into the development model and the larger [DevSecOps](http://blog.evident.io/blog/2015/3/26/the-marriage-of-devops-secops) movement. By default, Sparta functions execute with a minimal set of IAM privileges. These include **cloudformation:*** privileges so that Sparta lambda functions can [discover](http://gosparta.io/docs/discovery/), at lambda execution time, outputs of other infrastructure resources created by your application. The default privilege set, which is also appended to caller-defined _IAMRoleDefinition_ instances, includes:
`logs:CreateLogGroup
logs:CreateLogStream
logs:PutLogEvent
cloudwatch:PutMetricData
cloudformation:DescribeStacks
cloudformation:DescribeStackResource`

When available, Sparta ensures that the IAM Role policy document limits privileges to the minimal set of AWS Resources:




Instead of an _IAMRoleDefinition_ value, you can alternatively supply a string literal ARN value if you’re integrating Sparta into a broader ecosystem. Using ARN literals also eliminates the [CAPABILITY_IAM](http://docs.aws.amazon.com/AWSCloudFormation/latest/APIReference/API_CreateStack.html) requirement when new stacks are created.

#### CLI

Sparta-based applications are command line driven and use [goptions](https://github.com/voxelbrain/goptions) to parse and self-document the available options:
`$ go run main.go — help
Usage: application [global options] &lt;verb&gt; [verb options]``Global options:
 -n, — noop Dry-run behavior only (do not provision stack)
 -l, — level Log level [panic, fatal, error, warn, info, debug] (default: info)
 -h, — help Show this help``Verbs:
 delete:
 describe:
   -o, — out Output file for HTML description (*)
 execute:
   -p, — port Alternative port for HTTP binding (default=9999)
   -s, — signal Process ID to signal with SIGUSR2 once ready
 explore:
   -p, — port Alternative port for HTTP binding (default=9999)
 provision:
   -b, — s3Bucket S3 Bucket to use for Lambda source (*)`

#### Provisioning

The most common operation is provisioning: actually getting your Sparta Go function properly deployed to AWS Lambda. This workflow is defined in [provision.go](https://github.com/mweagle/Sparta/blob/master/provision.go) and we’ll walk through some of the major steps.**Preparation**

As part of supporting a specific application, Sparta needs a bit of plumbing to connect your functions to AWS Lambda. Most of this boilerplate is included during Sparta’s own build, which uses [mjibson](https://github.com/mjibson)/[esc](https://github.com/mjibson/esc) to [embed resources](https://github.com/mweagle/Sparta/blob/master/CONSTANTS.go) as part of its `go generate` phase (see [doc.go](https://github.com/mweagle/Sparta/blob/master/doc.go)). These embedded resources define the NodeJS layer used to connect Go to one of the supported [programming models](http://docs.aws.amazon.com/lambda/latest/dg/programming-model-v2.html). For each AWS Lambda addressable function, Sparta dynamically creates an exported Javascript function name using the _createForwarder_ function embedded in [resources/index.js](https://github.com/mweagle/Sparta/blob/master/resources/index.js):




The **path** argument in _createForwarder_ is extracted from your Go function via Go reflect:




If you include the ` — level debug` global command line option you can see the automatically generated JavaScript code (grep for _Dynamically generated NodeJS adapter_).**Compilation &amp; Generation**

After verifying your application’s IAM Roles, the next step is to build the application. This involves:

*   go generate: Running the [generate](https://blog.golang.org/generate) command against your codebase. This enables your application to do things like use [Hugo](http://gohugo.io/) to [build a static S3 site](https://github.com/mweagle/SpartaHugo).
*   Compile: Compile for Amazon Linux AMI



*   Code generation: Dynamically create the JavaScript code that exposes your Go functions (and Sparta’s own Lambda-backed CloudFormation Resources) as independently addressable AWS Lambda entities:



*   Package: package everything up into a [ZIP archive](http://docs.aws.amazon.com/lambda/latest/dg/nodejs-create-deployment-pkg.html), including the supporting node_modules that Sparta [embedded in its own build](https://github.com/mweagle/Sparta/blob/master/resources/provision/package.json).**Upload**

The next step is to upload the primary code archive, and if it’s defined, the supplementary [S3site](http://gosparta.io/docs/s3site/) zip archive that contains your site’s static resources. Not much interesting to see here, let’s move along.**Marshal &amp; Creation**

Sparta only uses CloudFormation to provision AWS resources. For areas where CloudFormation support lags behind available other AWS service features, Sparta uses [Lambda-backed CustomResources](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources-lambda.html). These were initially written in NodeJS, but starting with release 0.5.5 have begun to be migrated to Go via the [cloudformationresources](https://github.com/mweagle/cloudformationresources) package.

There are three primary Sparta types which are responsible for self-exporting their object graph to CloudFormation via [https://twitter.com/@crewjam](https://twitter.com/@crewjam) ‘s excellent [go-cloudformation](https://github.com/crewjam/go-cloudformation) package:

*   _sparta.LambdaAWSInfo_: This also handles marshaling all referenced types (Permissions, EventSources)
*   _sparta.API:_ All things APIGateway related
*   _sparta.S3Site:_ Everything related to provisioning, configuring, and populating an S3 bucket for static website hosting

Once the core template is assembled, it’s annotated via any user supplied [TemplateDecorators](https://godoc.org/github.com/mweagle/Sparta#TemplateDecorator) and ready for provisioning. Sparta checks the application state (using the **serviceName** discussed above) and either calls createStack or updates the existing one via ChangeSets.**Runtime**

At some point, even in this #serverless world, there is something (let’s just call it the æther), that is going to execute your AWS Lambda function. For Sparta applications, the registered AWS Lambda function name was defined as an _index.js_ export during the **Creation** phase, where a new _exports_ entry was written:




The exported name is the one registered with AWS Lambda and ultimately proxies the incoming AWS Lambda event to the sidecar Go application’s ServeHTTP function. ServeHTTP uses the HTTP path component as the lookup key into the Go function dispatch map and forwards if there’s a match:




The registered Sparta lambda function handles the request, the response to which is translated to an AWS Lambda-compatible result back in the calling NodeJS layer:



![image](/posts/2016-04-26_sparta-behind-the-shield/images/2.gif#layoutTextWidth)


That’s the core workflow, minus some gory details, that enables a Sparta-based application to self-provision to AWS Lambda. Sparta leverages both Go’s existing SDLC tools (go generate, [httptest](https://golang.org/pkg/net/http/httptest/)) and AWS Lambda’s flexibility to bring together security, operations, and application logic in a developer-centric manner. We’ll all be better off, both from engineering and business perspectives, the more we can break down the walls artificially separating those concerns and make it simple to securely, quickly, and profitably support the business.

![image](/posts/2016-04-26_sparta-behind-the-shield/images/3.png#layoutTextWidth)
Courtesy of [http://startupquote.com/post/5626579141](http://startupquote.com/post/5626579141)



### Future

In the next post I’ll look at how to add additional operational metrics and alert conditions to your Sparta application. Now that it’s provisioned, how can you tell what’s going on? And more importantly, be notified if it’s not going well?

In the meantime, if the **TODOs**, comments, or downright functionality can be improved, please open up a PR at [https://github.com/mweagle/Sparta](https://github.com/mweagle/Sparta). Ping me on Twitter @mweagle if you want to talk about Sparta or anything #serverless.
