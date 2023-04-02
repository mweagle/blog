---
title: "A Go framework for AWS Lambda"
author: "Matt Weagle"
date: 2016-04-11T18:15:22.125Z
lastmod: 2023-02-18T21:57:19-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "In previous posts I discussed why I think Lambda represents a new development opportunity. Lambda based services are operationally aware by…"

image: "/posts/2016-04-11_go-framework-for-aws-lambda/images/1.png"
images:
 - "/posts/2016-04-11_go-framework-for-aws-lambda/images/1.png"
 - "/posts/2016-04-11_go-framework-for-aws-lambda/images/2.jpeg"

aliases:
    - "/a-go-framework-for-aws-lambda-ab14f0c42cb"

---

![image](/posts/2016-04-11_go-framework-for-aws-lambda/images/1.png#layoutTextWidth)


In [previous posts](https://medium.com/@mweagle) I discussed why I think Lambda represents a new development opportunity. Lambda based services are operationally aware by design. Delegating some operational responsibility to AWS frees your team to focus on creating business value at higher velocity. Earlier PaSS offerings promised similar benefits, but those often didn’t materialize for a host of reasons. Lambda-based architectures provide a significantly different abstraction (FaSS — Functions as a Service?) that truly moves application development towards **#LessOps**.

This lofty pronouncement is nice and everything, but idle speculation is worth the pixels it costs. To promote the broader adoption of AWS Lambda (and [Go](https://golang.org/), for that matter), I’ve released [Sparta](http://gosparta.io): a (opinionated) Go framework for AWS Lambda microservices.

This post provides a brief introduction to Sparta, including primary design goals and a simple Sparta application. Sparta is written in Go, which AWS Lambda doesn’t [officially support](http://docs.aws.amazon.com/lambda/latest/dg/programming-model-v2.html). Isn’t AWS Lambda already a big enough shift? Why Go?

#### Go cloud, young people

> [](https://twitter.com/dberkholz/status/705279952128598017)


Basically, because Go offers a solid toolchain and set of primitives to write services. In no particular order, Go offers:

*   Single binary deployment
*   [Excellent concurrency primitives](https://blog.golang.org/pipelines)
*   An official [AWS Go SDK](https://aws.amazon.com/sdk-for-go/)
*   Extremely fast compilation
*   [Well-defined error handling patterns](http://blog.golang.org/error-handling-and-go)
*   Static types
*   Minimal startup overhead
*   [Cross-platform compilation](http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5)
*   [Rich standard library](https://golang.org/pkg/)
*   It’s [boring](http://stevebate.silvrback.com/go-is-boring), in the best possible way

As with all languages and tools there are some idiosyncrasies (I’m not a big fan of [_make_](https://golang.org/pkg/builtin/#make) &amp; miss _map/reduce_), but overall it’s a great fit for services.

It’s my hope that AWS Lambda will soon add first class support for Go (is [Tim Wagner](https://twitter.com/timallenwagner) in the house?), but in the meantime the AWS Lambda environment is rich enough to support Go binary execution.

#### HTTP All The Things

When I started Sparta, [Eric Hammond](https://alestic.com/2014/11/aws-lambda-environment/) and [Tom Maiaroto](https://serifandsemaphore.io/go-amazon-lambda-7e95a147cec8#.8b5myjmfa) had already shown it was possible for AWS Lambda to execute arbitrary executables included in a code archive. This meant given a little help from one of the officially supported runtimes, a compiled Go binary could be the target of an AWS Lambda request.

For Sparta, that help came in the form of a NodeJS HTTP proxying tier. For each Sparta-compatible AWS Lambda function (defined below), Sparta creates a unique NodeJS proxy route to forward the AWS request to a sidecar Go binary. The HTTP request/response semantics are a good fit for Lambda and allow applications to support one-time initialization as opposed to a process-per-request model (subject to [container reuse](https://aws.amazon.com/blogs/compute/container-reuse-in-lambda/)).

After a bit of work, the full HTTP-based proxying solution resulted in a workflow where Sparta:

*   Cross-compiles the Go binary for AWS’s Linux
*   Dynamically creates a NodeJS HTTP proxy entry for each unique Sparta lambda function. Each entry represents an addressable AWS Lambda function.
*   Builds a deployable ZIP archive that includes the Go binary and dynamically created JS contents.
*   Creates a CloudFormation template for provisioning, using content-based resource names
*   Either creates or updates CloudFormation stack

See the [Sparta docs](http://gosparta.io/docs/) for more details. In future posts I’ll dive more into the implementation as well as some limitations.

The end result of this is a 1:1 mapping of registered Go functions to their AWS Lambda counterparts. When an AWS Lambda function is triggered, the NodeJS shim ensures the Go process is available, and after receiving a SIGUSR2 handshake, manages the AWS Lambda interaction. Decoupling the AWS Lambda request from the Go function signature does incur some un/marshaling overhead, but I’ve not yet found that to be a blocker. On the positive side, using HTTP enables Sparta to use Go’s [httptest](https://golang.org/pkg/net/http/httptest/) for [localhost testing](http://gosparta.io/docs/local_testing/) .

Once the basic mechanics were in place, the next question was working out what a “Sparta” application would look like.

#### Going To Battle

One of Sparta’s primary design goals is to enable a secure, comprehensive, and self-contained specification of an AWS Lambda service. Security policies, supporting AWS infrastructure (eg, ElastiCache, DynamoDB, S3 buckets), logging, metrics, and alert triggers should be represented in a way comparable to business logic. After all, they’re all interrelated in production; shouldn’t that be reflected in the development model? There were other goals (using the compiler whenever possible, relying on CloudFormation for all provisioning operations, allowing arbitrary application structure/layout), but the primary goal remained.

It wasn’t realistic to represent all these dimensions before AWS Lambda. There were too many levels, and they each had their own wildly different models and requirements: [**rsyslog**](http://www.rsyslog.com/) and [**HAProxy**](http://www.haproxy.org/) configuration go together like peanut butter and jellyfish. Well, technically it’s software and could have been done, but it would have been DSL [astronauting](http://www.joelonsoftware.com/articles/fog0000000018.html) and definitely wouldn’t have ended well. However, with AWS Lambda concealing much of the operational complexity and AWS APIs available for the other primitives, the surface area is significantly reduced.

#### Hello World

With that background, let’s take a look at a Sparta application:




There’s a few items worth pointing out in this example:

**Sparta Lambda Functions**

The Sparta-compliant [LambdaFunction](https://godoc.org/github.com/mweagle/Sparta#LambdaFunction) signature is a bit different from the standard AWS Lambda ones. In addition to the AWS Lambda [context and event](http://docs.aws.amazon.com/lambda/latest/dg/nodejs-prog-model-handler.html), Sparta provides:

*   _http.ResponseWriter_: A writer to use for the AWS Lambda response data. Depending on the HTTP status code (success: 200–299), the response body is used as the AWS Lambda response or error data.
*   _*logger_: A [logrus](https://github.com/Sirupsen/logrus) instance that is preconfigured to produce JSON output to be consumed by CloudWatch Logs.

**Security Policies**

As part of registering a Sparta lambda function, callers must provide an IAM Role definition under which the AWS Lambda function will execute. IAM Roles are used to limit privileges to other AWS resources and prevent unwarranted breaches.

The example above provides an empty policy (_sparta.IAMRoleDefinition{}_) which prevents _helloWorld_ from accessing any other AWS resources. It’s also possible to define new, resource-targeted roles within an application and/or reference pre-existing Role ARNs. For example, see the [SpartaImager](https://github.com/mweagle/SpartaImager) service that uses a privately defined S3 bucket to store stamped images.

**Layout**

Sparta requires a Go application with a **main** package that ultimately calls _sparta.Main_ with the set of registered functions. No additional project structure is required.

#### To The Cloud

The result of this is a self contained Go application that can be deployed to AWS Lambda:
`&gt; go get -u ./...
&gt; go run application.go provision --s3Bucket {MY_S3_BUCKET}`

Where **{MY_S3_BUCKET}** is an S3 bucket that your currently configured [Go AWS SDK](http://docs.aws.amazon.com/sdk-for-go/latest/v1/developerguide/setting-up.title.html) credentials can access. After some log output:




You’ll see a new AWS Lambda function in the console:

![image](/posts/2016-04-11_go-framework-for-aws-lambda/images/2.jpeg#layoutTextWidth)


Which you can then invoke. When you’re done, delete the service via
`&gt; go run application.go delete`

#### No Lambda Is an Island

Sparta treats a Go package as a deployable entity. The package namespace denotes a logical grouping of interrelated functions which is deployed in a single operation. There is no a priori limit on the number of lambda functions defined by Sparta-based application. It’s very much possible to build a monolith with Sparta.

While this could be seen as a disadvantage, I think it’s part of the larger conversation around the need for new [serverless patterns](https://medium.com/@PaulDJohnston/we-need-more-serverless-patterns-17440704773a#.8nhgyi7zh). In my experience, services are small networks of interrelated calls, not only isolated functions.

For instance, consider a simple lambda function that mutates a DynamoDB record. To support search, the function should also update an ElasticSearch cluster. Handling both operations in a single function call opens up consistency and response latency issues. Treating them atomically creates potential schema and orchestration issues. I don’t know if “application” is the right term or concept, but I do think we’ll see new patterns emerge as more serverless services are created.

### Future

Hopefully this provided a bit of background on Sparta’s goals and the rationale behind choosing Go. I don’t know if or when AWS plans to officially support Go, o̶r̶ ̶w̶o̶r̶s̶e̶ ̶y̶e̶t̶,̶ ̶e̶l̶i̶m̶i̶n̶a̶t̶e̶ ̶t̶h̶e̶ ̶“̶b̶i̶n̶a̶r̶y̶ ̶l̶o̶o̶p̶h̶o̶l̶e̶,̶”̶ (**Correction**: running binaries is [officially supported](https://aws.amazon.com/blogs/compute/running-executables-in-aws-lambda/). Thanks [@applesaucefever](http://twitter.com/applesaucefever) ), but things seem to be stable for the time being. Over the next few posts I plan to cover:

*   Sparta internals
*   Metrics and alerting
*   Responding to other AWS Lambda [event sources](http://docs.aws.amazon.com/lambda/latest/dg/intro-core-components.html)
*   Provisioning &amp; discovering dynamic AWS infrastructure
*   API Gateway
*   Provisioning CORS-enabled static websites as part of an application

To learn more about Sparta please visit the [documentation](http://gosparta.io/docs/) and even better, [open a PR](https://github.com/mweagle/Sparta) to make things better. There are also several [sample applications](https://github.com/mweagle?tab=repositories) on GitHub. And please leave a comment or message me @mweagle if you’re running into Sparta issues or using it in production.
