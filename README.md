# Go Microservice Architecture

_This project is licensed under the MIT License._

## What & why

The intent of this project is to provide an example microservice that is scalable in terms of complexity. In order to achieve that, it has implements Domain-Driven Design patterns for separation of concerns. The domain core is located in the `api` directory, exposing some user-related logic. It explicitly does not import any non-domain package.

The microservice in this example exposes a graphql api over http, although the very purpose of this project is to design a microservice to be agnostic of how it receives its invocations of behaviour.

All technical implementation is initialised in the `main` function and injected as dependencies in their `API` interfaces, which allow the domain core to simply call a method that will figure out how to do. Hence, the domain core can focus on what should be done. It's easy to refactor implementation details, easy to test and more expressive.

## Getting started

In your deployment pipeline, build the service binary using the make target and provide as entrypoint in a containerized environment of your choice.

### Development

`cp sample.env local.env` and change accordingly.

In order to interate on this project, install the following:
- `go` (>=1.17)
- `golangci-lint`
- `gqlgen`

Have a postgres instance running and run the `scripts/db.sh` to create database and tables as configured in local.env.

## Evolutionary design

Another goal of this project is to showcase how a microservice can be set up to modularise sets of functionality and invoke behaviour in other components through an interface, without the hassle of having maintaining multiple services. One can choose a single service binary with multiple `API`s, located in `api/`. This can be useful in early stages of development, or when the time for breaking up a well-designed modularised monolith never comes. 

One can freely decide to promote a service's `api` when there are reasons to do so. See how the `email` api could easily be extrapolated as a service on its own. It acts as a broker between services to send messages to each other. This process should be as easy as moving the api to a separate process and replace the behaviour invocation implementation to do a network call or configure the event bus, depending on the communication pattern of the application. The intent is to have a clean architecture, such that the developer may easily promote an api to be its own microservice and just change the implementation of the interface - through which other apis would interact with it - to be a client's network call or event over the application's event store.

### Asynchronous vs Synchronous

Currently, the implementation of how service modules communicate is done via asynchronous events. Implementing a synchronous request-response type of communication is trivial - analogous to how a database client or any other network call may be implemented.

### WIP

This is still a work-in-progress. The packages serve their purposes, but may be renamed and restructured in a suitable hierarchy when one is deemed appropriate.

## Domain-driven-design
The core feature is to showcase a microservice architecture with an inward dependency direction of the following:

### Domain core (api/*)

In this example, `api/auth` and `api/email` contain the domain logic. It should use the ubiquitous language of the problem domain, and not include any technical implementation of any client call. 

All output (side-effects) other than a direct response back through the input layer proceeds via dependency-injected interfaces. 

1. The domain core defines these interfaces. 
2. The client layer implements them. 
3. The main function injects the clients as dependencies into the domain core, so it remains agnostic.

See the email api's `EmailClient` field as an example of how the domain core knows _what_ to do and _when_ to do it, but only the implementation of its interface in the client layer, injected as a dependency by `main`, knows _how_. In this case that is through a postmark client, but this can easily be substituted.

### Config layer
`config` has all the configuration that the microservice persists throughout its lifetime. It is strictly immutable. It includes configuration of the microservice's dependencies to be loosely coupled to its deployment environment. Only the `main` package may load the config from the inputs and passes only the relevant parts down.

- Environment variables: the obvious choice to configure a port, graphql endpoint, database configurations, etc.
- Command-line flags: if the microservice accepts command-line flags, they should be loaded in this package, included in the Config variable and returned to the main package.
- Filesystem data: files that are accessible at runtime can be loaded into the service's memory in this package. Mind that it is meant for immutable data. If your files are mutable and service-scoped (probably you don't want this), it should provide a different package to interact with. If it is shared across replicas of this microservice, it's seen as a storage implementation and hence be included therein.

### Input layer

The input layer is the interface through which other services invoke behaviour in this microservice. It transports application protocol-level information to domain core-level types and function calls. Two sources of input are servers and event subscribers. Both are implementations of processes that listen for messages, invoke the corresponding domain behaviour and output messages.

#### Server

`server` initialises an HTTP server that listens for requests, invokes domain behaviour and send a response. That behaviour might be limited to merely retrieving data from a repository. If one wishes, multiple servers can be described here that call the same domain logic and uses the same repository interface, but serve a different protocol, like REST, GraphQL, SOAP or gRPC.

#### Event subscriber

`eventsubscriber` starts a process that listens for events, invokes behaviour and publish zero or more new events.

Publish/subscribe, event-driven architectures can also be implemented by starting a similar process that takes as input the interface of a domain core; in this case an implementation of the auth.Auth interface. The server in that case is a 

### Client layer

`client` contains the implementations of any clients that the APIs need to perform output as side-effects. The interfaces are defined in the domain core, but they're implemented here, such that the domain core remains agnostic of any network dependencies, third-party libraries, storage interfaces, etc.. They are injected as dependencies in the service APIs upon initialisation.

#### Repository 

`repository` implements the client to the database of one or more entities or aggregates. The interface defined in the api, similar to other dependencies injected into the api. It translates the retrieval and persistence of entities that need to be committed to storage, and its functionality is called by the domain core package.

##### Storage 

`storage` resides in the repository and currently has a postgres database as its storage interface. The repository layer uses this database interface to implement the interface that the domain core needs to perform business logic.

