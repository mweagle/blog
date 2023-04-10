---
title: "See Lambda Go"
author: "Matt Weagle"
date: 2017-05-03T03:16:28.419Z
lastmod: 2023-02-18T21:57:31-08:00
tags: ["serverless", "sparta"]
description: ""

subtitle: "The Sparta framework for AWS Lambda provides a full-featured environment for writing Go-based, Serverless-powered applications."

image: "/posts/migrated/2017-05-03_see-lambda-go/images/1.png"
images:
 - "/posts/migrated/2017-05-03_see-lambda-go/images/1.png"
 - "/posts/migrated/2017-05-03_see-lambda-go/images/2.png"
 - "/posts/migrated/2017-05-03_see-lambda-go/images/3.png"
 - "/posts/migrated/2017-05-03_see-lambda-go/images/4.png"
 - "/posts/migrated/2017-05-03_see-lambda-go/images/5.png"
 - "/posts/migrated/2017-05-03_see-lambda-go/images/6.png"
 - "/posts/migrated/2017-05-03_see-lambda-go/images/7.png"
---

![image](/posts/migrated/2017-05-03_see-lambda-go/images/1.png#layoutTextWidth)


The [Sparta](http://gosparta.io) framework for AWS Lambda provides a full-featured environment for writing Go-based, Serverless-powered applications.

Up until version 0.11.1 however, there was a measurable overhead to this approach, especially in container cold-start scenarios. Cold-starts happen when a serverless function of any type needs to be created, placed, and invoked by a provider. For Sparta applications, this meant a relatively large performance penalty before the NodeJS-based shim was able to forward the request to the standalone Go binary.

Specifically:

1.  Unzipping the code bundle (handled by AWS)
2.  Copying the Go binary to _/tmp_
3.  Changing permissions and launching the binary
4.  Waiting for the Go binary to send a `SIGUSR2` signal to the parent NodeJS process, indicating it was fully initialized and ready to handle requests over HTTP

In addition, there was also measurable overhead for each request/response during marshaling/unmarshaling through JSON at the NodeJS layer. Take a look at the full [index.js](https://github.com/mweagle/Sparta/blob/master/resources/index.js) source for the full initialization and keep-alive behavior.

Sparta 0.11.1 introduces a new `cgo` target for existing applications that replaces the NodeJS shim with a [Python ctypes](https://docs.python.org/3/library/ctypes.html) interface that proxies your Sparta “application”. However, by the time your lambda is called, Sparta has already used the power of Go’s AST manipulation features to transform your application source into a library.

The end result is that by changing a single source import, you can achieve significantly better AWS Lambda performance.

![image](/posts/migrated/2017-05-03_see-lambda-go/images/2.png#layoutTextWidth)


No muss, no fuss. And most importantly, no changes to the rest of your Sparta application.

The following sections walk through the three major parts of this feature, using the `cgo` application in the [SpartaPython](https://github.com/mweagle/SpartaPython) project.

![image](/posts/migrated/2017-05-03_see-lambda-go/images/3.png#layoutTextWidth)


### Go Rewriting

I’m generally very reluctant to rewrite user input. Much of the time it’s sensible to treat user input as immutable and err on the side of validation and actionable error messages if some precondition isn’t met.

That didn’t seem to be an option here though, as I definitely didn’t want the downstream target runtime to pollute the general idea of creating a AWS Lambda function via `sparta.NewLambda(…)`. Introducing something akin to `sparta.NewLambdaForPythonWithBetterPerformance(…) `seemed like the wrong approach.

After exploring several other options, I concluded that source rewriting was the only viable option, and started to build up the transformations. What I ended up with is a source rewriting pipeline that does the following:

1.  Parses the application input: The new `spartaCGO.Main(…)` call site file MUST be your application’s `package main` source file. Sparta [locates the source input](https://github.com/mweagle/Sparta/blob/master/cgo/cgo_main.go#L22) so that it can make the application’s initialization logic library compatible (see below).
2.  Adds `import` statements: This uses the [golang.org/x/tools/go/ast/astutil](https://godoc.org/golang.org/x/tools/go/ast/astutil) package which provides higher-level manipulation functions. While it’s possible to drop down a level to rewrite the source, I was very grateful for the one-liners.
3.  Renames your application’s`main` function to `init` so that your functions are properly registered when called via a library context.
4.  Adding the `cgo`-exported [function](https://github.com/mweagle/Sparta/blob/master/cgo/walkers.go#L20) that will be called by the Python AWS lambda shim

The end result of this is modified version of your source that defines a `cgo` build target. Once your original source is renamed to a non-Go source file in the input directory, the transformed source is exported and ready for building.

#### Before




#### After




Note that the transformed file is also preserved in the _./.sparta_ build working directory as part of a build. If you provision with the `--noop` flag, all build artifacts including the transformed source are available in that same directory.

![image](/posts/migrated/2017-05-03_see-lambda-go/images/4.png#layoutTextWidth)


### Docker

The next step is to build the source, including the `--buildmode=c-shared `flag as part of the command line arguments. The build itself is done in Docker, using a volume mount mapping so that the library binary can be referenced afterwards:

`INFO[0000] Building `cgo` library in Docker Args=[run — rm -v /Users/mweagle/Documents/gopath:/usr/src/gopath -w /usr/src/gopath/src/github.com/mweagle/SpartaPython/cgo -e GOPATH=/usr/src/gopath -e GOOS=linux -e GOARCH=amd64 golang:1.8.1 go build -o SpartaHelloPythonCGO.lambda.so -tags lambdabinary linux -buildmode=c-shared -tags lambdabinary noop ] Name=SpartaHelloPythonCGO.lambda.so`

Sparta will automatically include any `SPARTA_` prefixed environment variables to the Docker build command in the event that the service requires host-level configuration. See the [cgo](http://gosparta.io/docs/cgo/) documentation for more details.

![image](/posts/migrated/2017-05-03_see-lambda-go/images/5.png#layoutTextWidth)


### Python

The final step is to create the Python shim that connects the AWS Lambda functions to their Sparta implementations through the exported `Lambda` function. The Python forwarding logic is relatively generic and is responsible for assembling the C-library analog to the NodeJS HTTP-based version.

The request is supplied to the `cgo` exported `Lambda` function which Sparta injected into the source:




The Go Lambda exported function handles converting the C types to their Go counterparts and updating the AWS credential information that’s available to the Python runtime. Once everything is initialized, Sparta ultimately makes an “HTTP request” to your function.

This isn’t strictly accurate though, since what Sparta actually does is create a synthetic, in-process mock request that exploits the existing HTTP-based Sparta lambda function signature:

`func HelloWorld(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger)`

The library version leverages this signature to make a “request” to the HTTP dispatcher that maintains a map of your registered functions .




After the internal request is routed and handled by your lambda function, there’s some more marshaling data to `ctypes` so that they can be sent back to Python. Finally, Sparta sends off some internal metrics to CloudWatch about function performance.

For the complete details, please take a look at the [source](https://github.com/mweagle/Sparta/blob/master/cgo/cgo_main_awslib.go#L130).

![image](/posts/migrated/2017-05-03_see-lambda-go/images/6.png#layoutTextWidth)


### Activating

Sparta clients can take advantage of this functionality by changing an import statement and updating a function call site.

Add:




to your _main.go_ source and change your primary Sparta entry point from `sparta.Main(…)` to `spartaCGO.Main(…).` That’s it.

Because [cgo is not go](https://dave.cheney.net/2016/01/18/cgo-is-not-go), this is an opt-in feature. The existing Sparta interfaces and runtime path remains unchanged via `sparta.Main(...).`

### Results

While benchmarks are notoriously variable, I’ve seen coldstart times drop from the ~1500ms range to the ~500ms range on 128MB memory allocation instances.

There is also a significant per-call performance improvement as noted in the image above.

![image](/posts/migrated/2017-05-03_see-lambda-go/images/7.png#layoutTextWidth)


### See Lambda Go

This was a relatively complex feature to implement, primarily because I didn’t want to introduce API-incompatible changes and it was my first time using Go’s AST packages. And with great power, comes a great deal of time spent reading documentation.

Performance is a feature, and especially with Serverless pricing models, improved performance directly translates to the bottom line. While there are tradeoffs to both NodeJS and Python shims, you should be able to shift between them as your technical and business needs evolve.

If you’re already a Sparta user, please give the `cgo` entrypoint a try and confirm that you see similar performance improvements. And if you’re not a Sparta user, this new feature helps to bring Go performance in line with officially supported runtime performance. Start today!

#### Notes

*   Thanks to [ashleymcnamara](https://github.com/ashleymcnamara) for the [https://github.com/ashleymcnamara/gophers](https://github.com/ashleymcnamara/gophers) images
*   Other images via Google, labeled for non-commercial use
*   Questions, comments? Please open a [GitHub issue](https://github.com/mweagle/Sparta/issues) or get in touch with @mweagle at the [Serverless Slack Forum](http://wt-serverless-seattle.run.webtask.io/serverless-forum-signup?webtask_no_cache=1).
