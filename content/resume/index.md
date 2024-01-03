---
title: "Matt Weagle - Resume"
subtitle: "🌎 Seattle, WA"
draft: false
description:
date: 2024-01-03T12:13:06-08:00
lastmod: 2024-01-03T12:13:06-08:00
tags: ["resume"]
ShowToc: true
---

## Contact Me

Please visit the [homepage](http://mweagle.io) for contact information.

## Summary

Customer-focused, hands-on software engineering manager and individual contributor with experience managing technical teams and delivering products and platforms in fast-evolving consumer services spaces. Experience building highly available, fault-tolerant, large-scale distributed systems. Brings system-level thinking to design discussions. Servant-leader, pragmatic decision-maker, focused on sustainable delivery, and advocate for generative culture.

---

## Experience

### Software Development Manager - AWS (2020 - 2023)

Software Development Manager for [AWS Lambda](https://aws.amazon.com/lambda/) and [AWS CodeCatalyst](https://codecatalyst.aws/explore).

* Contributed to [AWS Lambda SmartStart](https://docs.aws.amazon.com/lambda/latest/dg/snapstart.html) and [AWS Lambda response streaming](https://aws.amazon.com/blogs/compute/introducing-aws-lambda-response-streaming/) capabilities.
* Defined and drove meetings, mechanisms and metrics to coordinate CI/CD for 100+ person development team spanning multiple organizations. The result of this work was increased visibility into per-component change/fail rates, regional deployment skew, and batch size. Metrics were automatically published as part of CI/CD execution and included standardized release note generation with commit history to reduce operational response time.
* Over the course of a year, doubled team size to increase capacity and reduce operational load.
* Managed decomposition of legacy `go` codebase into multiple, localhost `gRPC` microservices to align ownerhip, increase organizational velocity, and reduce response latency.
* Led the [Lambda Corection of Error Bar Raiser](https://wa.aws.amazon.com/wat.concept.coe.en.html) program to improve our incident documentation, learnings, and discussions.
* Delivered [GitHub Actions](https://aws.amazon.com/blogs/devops/using-github-actions-with-amazon-codecatalyst/) on CodeCatalyst.
* Owned delivery of [CodeCatalyst Build](https://docs.aws.amazon.com/codecatalyst/latest/userguide/build-workflow-actions.html) and [CodeCatalyst Workflow Actions](https://docs.aws.amazon.com/codecatalyst/latest/userguide/workflows-actions.html).

---

### Engineering Manager - [Lyft](https://www.lyft.com/) (2019-2020)

Engineering Manager for Forecasting and ML Training/[Flyte](https://flyte.org/) teams. The Lyft Forecasting team is responsible for providing near real-time (~5 minute) and long-term business forecasts (~12 month).

* Created space for the team to replace legacy system with Spark-based pipelines for scalable, reliable, and cost-efficient training. Worked with team to migrate a near real-time system to a more durable and scalable solution using Kafka, [Flink](https://flink.apache.org/), [Druid](https://druid.apache.org/) and [Seldon](https://www.seldon.io/).
* Fostered a culture of learning, collaboration, and delivery.
* Mentored and coached within and across the team.
* Instituted [press release driven development](https://www.allthingsdistributed.com/2006/11/working_backwards.html) practices to clarify objectives and closeaccountability gaps.
* Conducted bi-weekly retrospectives to recognize accomplishments and discuss opportunitiesfor improvement.
* Cross-functionally collaborated with ML scientists to reduce downstream engineeringintegration costs.

The Flyte team (now [Union.AI](https://www.union.ai)) was responsible for providing a platform for 100,000+ workflow executions a month for business-critical workloads.

* Maximized team retention during organizational restructuring.
* Changed on-call and escalation policy to improve user experience, reduce on-call burden, and increase system feedback cycles.
* Prioritized observability work to facilitate more customer-facing transparency which allowed for greater self-service.
* Identified opportunities for systemic reliability improvements.

---

### Director of Site Reliability Engineering - [ShiftLeft](https://shiftleft.io/) (2016-2019)

Responsible for bootstrapping core *go* codebase for common microservice architecture which represents a processing graph mediated by Kafka topics. Successfully onboarded additional team members to evolve and extend the core codebase. Additional accomplishments:

* Ensured all services expose Prometheus-compatible observability information.
* Created single fail-fast configuration initialization path with actionable error feedback.
* Developed integrations with Kafka, Kinesis, Kinesis Firehose, Minio, S3, Prometheus, Segment, GitHub, Vault, and other third-party services.
* Built superuser GUI in ReactJS to support authenticated BackOffice style operations and secure access to encrypted-at-rest customer assets.
* Initiated shared documentation build pipeline and infrastructure.
* Created shared, tagged, and structured logging interface.
* Collaborated with operations to ensure successful cluster-agnostic feasibility.
* Created initial E2E integration pipeline to ensure ongoing service compliance.
* Provided ongoing operational support.

---

### Engineering Manager - [NodeSource](https://nodesource.com/) (2016-2016)

Managed and contributed to team responsible for building [NodeSource Certified Modules](https://docs.nodesource.com/ncm/docs), an SaaS offering built on AWS using EC2, SQS, CloudFormation, and several other AWS service offerings. Highlights of my work:

* Created JavaScript toolchain to automate CloudFormation service provisioning.
* Worked with sales team to provision NodeSource VMs to AWS Marketplace.
* Delivered CloudFormation-backed build pipeline that provisioned ephemeral EC2 instances to build multi-cloud VMs with [Packer](https://packer.io/).

---

### Engineering Manager - Adobe (2014-2015)

Led engineering team to deliver the [Adobe Creative Cloud Assets](https://assets.adobe.com/) presence. The team is responsible for building and supporting all service aspects from the DevOps/infrastructure layer to the business-logic microservices tier. All services are deployed to AWS across multiple accounts and geographies. In addition to supporting large-scale release events, helped lead a feature-flipper gated system that enabled continual product improvement. Additional accomplishments:

* Prioritized and managed the team’s technical roadmap.
* Mentored and grew the team.
* Migrated the platform services tier to Docker and etcd-backed configuration service. Migration occurred in tandem with CI/CD pipeline delivering multiple releases per day.
* Identified build-vs-buy opportunities and secured funding as needed.
* Provided direct operational support for all developed services (i.e.: “You build it, you run it”).
* Established lightweight processes to improve delivery cadence and quality.

---

### Engineering Lead - Adobe (2012-2014)

Lead architect and developer for Adobe Creative Cloud [Assets](https://assets.adobe.com/) services tier. Highlights of my work:

* Led complete application rewrite (from Ruby to Node.js) and MongoDB database denormalization to alleviate performance bottlenecks. Twice lazily migrated millions of user datasets with zero downtime.
* Defined Node.js coding standards & architectural patterns.
* Designed restify-backed core workflow engine to minimize response latency during HTTP request handling.
* Integrated additional AWS services for improved performance and resiliency.

---

### Senior Computer Scientist - Adobe (2010-2012)

Developer of real-time, collaborative [video editing product](http://www.adobe.com/products/adobeanywhere.html). The server is delivered as a set of OSGi bundles installed into custom distributions of Apache Sling and Adobe CQ5. These runtimes embed the Java Content Repository and provide a NoSQL-style hierarchical storage model together with a Sling HTTP routing mechanism. Accomplishments:

* Designed and implemented JVM-based job handler. The agent is responsible for consuming JMS messages and delivering them to a locally managed native process.
* Implemented logfile aggregation, arbitrary process monitoring, and host metric sampling. Exposed observability information via the SIGAR library, published as JSON over-HTTP, and visualized with Dygraphs.
* Implemented in Scala using Akka and Apache Camel.

---

### Senior Computer Scientist (Lead) - Adobe (2007-2010)

Engineering lead for Atom Publishing Protocol solution for rich-content delivery. This end-to end solution featured a rich Internet application media upload tool for batched, localizable uploads to an AtomPub server. Highlights of my work:

* Implemented REST-compliant API built on a stack of Java/Restlet, XQuery, and Lucene.
* Designed and implemented multiple memory caching strategies and render-to-disk phase in the publish-to-production step to minimize latency.
* Created RIA end-user client that dynamically reconfigured UI in response to server state. Client included inline video playback, dynamic view theming, and native drag-and-drop functionality.

---

### Computer Scientist (Lead) - Adobe (2004-2007)

Project lead for cross-platform installer package compiler that uses an XML input file with declarative semantics to build Windows Installer Packages and proprietary Macintosh OSX archives. Accomplishments:

* Technology was adopted by more than 20 internal teams within a single 18-month product release cycle.
* Created the system that uses XML and Regular Expressions to allow users to inductively declare install-time behaviors. Build-pipeline is composed of Python, XSLT, C++, and C#/Cocoa. Cross-platform runtime user interface uses embedded Internet Explorer control supported by core DOM published via COM.

---

### Computer Scientist - Adobe (1999-2004)

Supported more than a dozen different Windows application installers. Responsibilities included:

* Requirements gathering
* Primary implementation
* Continuous delivery/integration
* Enterprise support
* Knowledge transfer
* Localization/internationalization

## Teaching Assistant - University of Notre Dame (1995-1997)

Prepared and delivered lectures designed to supplement core course material. Provided weekly lecture to three different student groups. Received consistently high positive teaching reviews from students.

## Community

* Creator of [Sparta](https://gosparta.io/) framework for *go* microservices powered by AWS Lambda
* [ServerlessDays Seattle Co-Organizer](http://seattle.serverlessdays.io/)
* [Serverless Seattle Meetup Organizer](https://www.meetup.com/Serverless-Seattle/)
* [Former AWS Serverless Community Hero](https://aws.amazon.com/developer/community/heroes)

## Education

### Certificate Program - University of Washington (2001-2002)

Completed with distinction the University of Washington program in Algorithms and Data Structures. This is a professional degree program in data structures, algorithm complexity, and discrete mathematics.

### Ph.D. Candidate - University of Notre Dame (1993-1997)

Doctoral candidate in Economics with focus on History of Political Economy and Computational Economics. Thesis work referenced three interrelated subjects: the history of artificial intelligence and economics researchers, the use of genetic algorithms for nonlinear optimization in economics and complex systems, and a theoretical model of market institutions based on computability theory. Presented research work on the history of artificial intelligence and economics at American Economic Association 1996 annual meeting.

### Summer Workshop - Santa Fe Institute (1996)

Attended [Computational Economics Workshop](https://www.santafe.edu/engage/learn/schools/graduate-workshop-computational-social-science-mod) on the role and importance of nonlinear optimization and complexity theory in economics. Coded genetic algorithm to explore evolutionary optimization concepts applied to market behavior.

### BA - Bucknell University (1993)

Majored in Economics with minor in English. Graduated *cum laude* and *Phi Beta Kappa* member.

## Patents

* [8768924](https://patents.justia.com/patent/9288248) - Conflict Resolution in a Media Editing System