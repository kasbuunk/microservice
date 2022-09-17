# Go Microservice Architecture

_This project is licensed under the MIT License._

## What & why

The intent of this project is to provide an example microservice that is _scalable in terms of complexity_. In order to achieve that, it applies Domain-Driven Design patterns for separation of concerns. Ports and adapters further help to isolate the domain core. The example exposes a graphql api over http, although the very purpose of this project is to showcase a microservice that is agnostic of the source that invokes its behaviour.

## Principles

- input: there are three sources of input, _requests_, _events_ and _commands_. Any and all requests are _served_ through a `server.Server` interface. Any and all events are _handled_ by the `event.Handler` interface. Commands have yet to be implemented, but follow the same isolation of input source from the invoked behaviour. 
- processing: the `app` (sub)packages do not import any other package. They determine all domain-specific behaviour, but delegate any side-effects to the implementation of client interfaces it defines under the `app/client` package.
- output: there are two destinations of output, _responses_ and _effects_. Reponses are the return value of a `server`'s request. _Effects_ (or side-effects) are any state change performed by a `client`, including: newly published events, requests out- or inside the application cluster and database queries through a `repository` interface.

All technical implementation is initialised in the `main` function and injected as dependencies in their `App` interfaces, which allow the domain core to simply call a method that will figure out how to do. Hence, the domain core can focus on what should be done. It's easy to refactor implementation details, easy to test and more expressive.

### WIP

This is still a work-in-progress. The packages' purpose and their interdependencies are a proof of concept. In the future, they may be renamed, merged or segregated and restructured in a suitable hierarchy. 

## Evolutionary design

A goal of this project is to showcase how a microservice can be set up to modularise sets of functionality and invoke behaviour in other components through an interface, without the hassle of maintaining multiple services. One can choose a single service binary with multiple `App`s, located in `app/`. This can be useful in early stages of development, or when the time for breaking up a well-designed modularised monolith never comes.

The intent is to have a clean architecture, such that the developer may easily promote an app to be its own microservice and just change the implementation of the interface - through which other apps would interact with it - to be a client's network call or event over the application's event store. See how the `email` app could easily be extrapolated as a standalone microservice. This process should be as easy as moving the app to a separate process and replace the behaviour invocation implementation to do a network call with a server in between, or configure the event bus, depending on the application's communication pattern. 

## Installation & development

In your deployment pipeline, run `make` to build the service binary `bin/app` and provide it as entrypoint in an environment of your choice.

`cp sample.env local.env` and change accordingly.

In order to iterate on this project, install the following:
- `go` (>=1.17)
- `golangci-lint`
- `gqlgen`

Have a postgres instance running and run the `scripts/db.sh` to create database and tables as configured in `local.env`.

### Asynchronous vs. synchronous messages

Currently, the implementation of how service modules communicate is done via asynchronous events. Implementing a synchronous request-response type of communication is trivial - analogous to how a database client or any other network call may be implemented.

## Separation of concerns

The core feature is to showcase a microservice architecture with an inward dependency direction of the following:

### Domain core (app/*)

In this example, `app/auth` and `app/email` contain the domain logic. It should use the ubiquitous language of the problem domain, and not include any technical implementation of any client call. 

All output (side effects) other than a direct response back through the input layer proceeds via dependency-injected adapters. 

1. The domain core defines its `ports` as interfaces. 
2. The `adapters` implement them. 
3. The main function initialises the `adapters` and injects them as implementations of the `ports` into the domain core. Thus the domain core remains agnostic.

See the email app's `EmailClient` field as an example of how the domain core knows _what_ to do and _when_ to do it, but only the implementation of its interface in the side effect layer, injected as a dependency by `main`, knows _how_. In this case that is through a postmark client, but this can easily be substituted.

### Config layer
`config` has all the configuration that the microservice persists throughout its lifetime. It is strictly immutable. It includes configuration of the microservice's `ports`, i.e. its dependency interfaces, to be loosely coupled to its deployment environment. Only the `main` package may load the config from the inputs and passes only the relevant parts down.

- Environment variables: the obvious choice to configure the port the process listeb on, graphql endpoint, database configurations, etc.
- Command-line flags: if the microservice accepts command-line flags, they should be loaded in this package, included in the Config variable and returned to the main package.
- Filesystem data: files that are accessible at runtime can be loaded into the service's memory in this package. Mind that it is meant for immutable data. If your files are mutable and service-scoped (probably you don't want this), it should provide a different package to interact with. If it is shared across replicas of this microservice, it's seen as a storage implementation and hence be included therein.

### Input layer

The input layer is the interface through which other services invoke behaviour in this microservice. It transports application protocol-level information to domain core-level types and function calls. Two sources of input are servers and event subscribers. Both are implementations of processes that listen for messages, invoke the corresponding domain behaviour and output messages.

#### Server

`server` initialises an HTTP server that listens for requests, invokes domain behaviour and send a response. That behaviour might be limited to merely retrieving data from a repository. If one wishes, multiple servers can be described here that call the same domain logic and uses the same repository interface, but serve a different protocol, like REST, GraphQL, SOAP or gRPC.

#### Event handler

`event.Handler` starts a process that listens for events, invokes behaviour. The behaviour may dictate to perform some side-effects and publish new events.

Currently, the event client is implemented to be in-memory, but the aim for it is to be an abstraction of how an eventbus client would interact with the application event store.

#### Command 

_TODO: Invoking app calls through a command has yet to be implemented._

### Adapter layer

`adapter` contains the implementations of any `port` dependencies the Apps need to perform output as side-effects. The interfaces are defined in the domain core as `ports`, but they're implemented here, such that the domain core remains agnostic of any network dependencies, third-party libraries, storage interfaces, etc. They are injected as dependencies in the service `Apps` upon initialisation.

#### Repository 

`repository` implements the client to the database of one or more entities or aggregates. The interface defined in the app, similar to other dependencies injected into the app. It translates the retrieval and persistence of entities that need to be committed to storage, and its functionality is called by the domain core package.

##### Storage 

`storage` resides in the repository and currently has a postgres database as its storage interface. The repository layer uses this database interface to implement the interface that the domain core depends on to perform business logic.

