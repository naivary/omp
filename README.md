# Cloud Native API Layout (cnapi)

This repository contains a set of best practices from the web and foundations
like CNCF providing a starting point to develop REST APIs in Go. To get started
just clone the repository and make changes to it as it would be your project.

## Environment

Cloud native applications are always delivered as container images and deployed
on Kubernetes (or some other Kubernetes flavor) for the orchestration of the
workloads. Given these circumntances it is save to say that Kubernetes is the
runtime of the Cloud and applications should be tailored to excel in these
conditions.

Furthermore, it is highly recommended for your dev environment to be identical
to your production environment to be more confident in your development and
tests. Great tools for local Kubernetes environments are `minikube` and `k3s`.
`kind` is anohter option but I woudln't recommended it because not all
functionalities are fully available especially in the storage realm.

## Stateless

Applications in containers execl when they are stateless. Meaning no data is
associate with the container or the application itself or in other words the
container could be deleted and started again without any effects on the
application. One example for stateful applications are databases.

Stateless applications are especially important because scaling horizontally is
straightforward by creating more replicas of the application.

## Service Discovery

If you application need to communicate with another service to run successfully
you need to somehow connect to the application. Using the IP-Address of the
container hosting the service in a Cloud-native environment is not recommended
because IP-Addresses are considered volatile and might change at any time. But
also it is not the responsiblity of the application to implement service
discovery.

The recommended approach to this problem is using your environment efficiently.
In Kubernetes this is done by creating an `Service` object. Applications which
need to connect to that service will be provided with the domain specific
connection string using the configuration options provided by the application
(see: [Config](#configuration))

## Configuration

Configuration is an essential part of every application. This makes an
application adaptable to different environments like staging, production or
development. Applications store config in different forms. It is highly
recommended to seperate Code from config. A litus test to check if an
application is correctly separting code from config is where the codebase could
be made open-source at any moment and no credentials would be exposed.

The 12FA App recommends to store configuration in the environment of the
application e.g. as environemnt variables. This allows for an language agnostic
approach without accidentaly checking configuration files into repositories.
Antoher option is using CLI flags which follow the principles of environment
variables. Both of the apparoaches have the advantage of adapting to the Cloud
environment e.g. Kubernetes because the values can be provided using various
Kubernetes sources (ConfigMap, Secrets etc.).

In Go you can easily implement both approaches. If you want to implement a CLI
use the `cmd/` directory for it. You can either do it by using the standard
library or using something like `cobra`. Environment variables can be retrieved
using `os.Getenv`.

For a dynamic approach in the code base the `run` function accepts a `getenv`
function which controls the config of the environment allowing you to configure
your environment as you wish.

## Principles of Chaos Engineering

Errors are part of every software and their occurence should be accepted as a
given fact. Because of that it is important to design APIs which are able to
cope with errors.

Chaos Engineering is enforcing exactly that. It is the discipline of
experimenting on a system in order to build confidence in the system's
capabiltiy to withstand turbulent conditions in production. For a detail
description consult the
[official documentation](https://principlesofchaos.org/).

To run local chaos experiments you can use
[Chaos Monkey](https://github.com/Netflix/chaosmonkey). A tool implementing the
principles of chaos engineering by Netflix. For chaos experiments on Kubernetes
use [Chaos Mesh](https://chaos-mesh.org/), a CNCF Incubating project. It is
recommened to choose the tool based on which is representing your production
environment the most.

## Documentation

Documentation of your application should be done in a _as code_ manner. Write
the documentation in markdown or mdx to allow for easy processing and generation
of content like a documentations website. The documentation of the API resources
and endpoints are implemented using the OpenAPI Standard. To reduce the amount
of boilerplate work you can use a tool which is generating the OpenAPI
Artificats based on your code. One example is
[codemark](https://github.com/naivary/codemark) which is generating any kind of
artifact based on (comment) markers.

## Validation

Validation is a crucial part of any HTTP request. Users can provide malicous
content or mismatch information. Therefor it is important to validate the
integrity of the payload before processing the request.

One easy approach to validate data is using a common interface for it. This
project provides the `Validator` interface in `validator.go`. The validator
takes in a context and provides feedback of the validations in form of a map. If
the returned map is empty then the validation was successfull. Combining this
approach with a JSON Schema is allowing for the most flexibile and industry
standard based solution.

## Testing

Testing is a crucial part of software development. Without it we cannot have any
confident in the funtionality of the application and development will be
staggered because of fear of regression. Testing can be done with different
goals is mind. For example unit tests assure that the smallest unit is working,
End-to-End (E2E) assure that every component involved in the request is
functioning correctly or load test assure that the application can scale
correclty.

Using the concept of Test driven development (TDD) makes it possible to
implement Unit Tests while implementing needed funtionality for the API.

E2E are a bit trickier because we need to wait for the server and it's
dependencies to be ready for incoming requests. For that the `probe` package can
be used. The package allows you to wait until the server is ready using public
faced API endpoints like `/readyz` or `/livez` to check for readiness.

E2E should be prioritised over Unit Tests because they are closer to the end
users experience. That might mean that Unit Tests implemented in TDD might be
deleted afterwards in favor of E2E tests.

## Telemetry

Exporting telemetry data of your application is highly important to get insight
into your application during production usage. Without it debugging, future
development and optimazation will be hard to accomplish. For an API three
telemetry kinds are relevant: Logs, Metrics and Traces. Starting with logs.

A Log is a timestamped text record, either structured (recommended) or
unstructured, with optional metadata. Go has a structured logging in the
standard library but other popular options exists like `zap`. Independent of the
library used to log it is important to log any data to `stdout/stderr`. Log
management is not the responsiblity of the application. Especially in a
Kubernetes environment there MUST exist a logging solution collecting relevant
logs produced by the pods.

Metrics are numerical measurements in layperson terms. Metrics play an important
role in understanding why your application is working in a certain way. For
example the amount of request over a defined period of time. The de-facto
standard for metrics format and collection in cloud native environments is
Prometheus a CNCF graduated project. Therefor it is important to make sure that
the library used for instrumentation is Prometheus compatible.

Traces give us the big picture of what happens when a request is made to an
application. Whether your application is a monolith with a single database or a
sophisticated mesh of services, traces are essential to understanding the full
“path” a request takes in your application. It is even more important in
distributed systems to be able to debug errors efficiently.

For all three telemetries it's highly recommened to use OpenTelemetry to
instrument the application. OpenTelemetry is a graduated CNCF project allowing
for a standardized way to instrument applications for all three telemetrie
kinds.

## Authentication and Authorization

Implementing authentication and authorization is not the responsiblity of the
service requiring it. There MUST be a dedicated central service for all other
service to be able to use for authentication and authorization based on standard
hardend protocols. These are mainly OpenID Connect and OAuth 2.0. Especially in
combination with an [API-Gateway](#api-gateway) authentication and authorization
SHOULD be handled centrally during ingress for all services.

If no central API Gateway like Kong is used and your application requires
authentication and authorization its recommended to use
[coreos/go-oidc](https://github.com/coreos/go-oidc) and
[golang.org/oauth2](https://pkg.go.dev/golang.org/x/oauth2) SHOULD be used.

## API Gateway

An API Gateway acts as a mediator between client applications and the backend
services within the microservices architecture. It is a software layer that
functions as a single endpoint for various APIs performing tasks such as request
composition, routing, and protocol translation. The API gateway controls
requests and responses by managing the traffic of APIs while enforcing security
policies. This simplifies API management by providing one central point of
control which aids developers in focusing on building individual services rather
than being encumbered by complex networks of APIs, including tasks such as user
authentication and rate limiting (see:
[Kong](https://konghq.com/blog/learning-center/what-is-an-api-gateway))

Using an API Gateway is enabling you to centralized many common functionalities
across all your services. In Kubernetes it is done by using the built-in kind
`APIGateway` in combination with an ingress controller.

Caching, TLS termination etc.

## Style Guide

It's important to have a unified style in your ever growing codebase. If more
people join they have their own ideas of how things should look like and it
makes it harder to mantain and extend your code. An unified experience can be
craeted by using style guides. These can be already existing internal ones or
popular public ones like
[Uber](https://github.com/uber-go/guide/blob/master/style.md) or
[Google](https://google.github.io/styleguide/go/).

Personally I like the Uber styleguide because it's easier to understand and
provides great examples to go through but the decision is upto you and your
company.
