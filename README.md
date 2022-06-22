# Go Microservice Architecture

This project is licensed under the MIT License.

The microservice in this example exposes a graphql api over http, although the very purpose of this project is to design a microservice to be agnostic of how it receives its invocations of behaviour. 

## What & why

The intent of this project is to provide an example microservice that is scalable in terms of complexity. In order to achieve that, it has implements Domain-Driven Design patterns for separation of concerns. The domain core is located in the `api` directory, exposing some user-related logic. It explicitly does not import any non-domain package.

All technical implementation is initialised in the main entrypoint and injected as dependencies in their `API` interface, such as an `EmailCLient` that allows the domain core to simply call a method that will figure out how to do. Hence, the domain core can focus on what should be done. It's easy to refactor implementation details, easy to test and more expressive.

## Evolutionary design

Another goal of this project is to showcase how a microservice can be set up to modularise sets of functionality and invoke behaviour in other components through an interface, without the hassle of having maintaining multiple services. One can choose a single service binary with multiple `API`s, located in `api/`. 

This can be useful in early stages of development, or when the time for breaking up a well-designed modularised monolith never comes. One can freely decide to extract a service when there are reasons to do so. 

This process should be as easy as moving the api to a separate project and replace the behaviour invocation implementation to do a network call or configure the event bus, depending on the communication pattern of the application.

See how the `email` service could easily be a service on its own, but is abstracted away by a general `MessageBus` interface. It acts as a broker between services to send messages to each other.

### Asynchronous vs Synchronous

Currently, the implementation of how service modules communicate is done via asynchronous events. Implementing a synchronous request-response type of communication will be done in a future iteration.

## Prerequisites

Install the following:
- `go` (>=1.17)
- `golangci-lint`

If you wish to further develop and change the schema, install `gqlgen`.

`cp sample.env local.env` and change accordingly.

Have a postgres instance running and run the `scripts/db.sh` to create database and tables as configured in local.env, and a `_test` database.

## How to use this microservice

In your deployment pipeline, build the service binary using the make target and provide as entrypoint in a containerized environment of your choice.

### WIP

This is still a work-in-progress. The packages serve their purposes, but may be renamed and restructured in a suitable hierarchy when one is deemed appropriate.

## Domain-driven-design
The core feature is to showcase a microservice architecture with an inward dependency direction of the following:

### Domain core (api/*)

In this example, `auth` and `email` contain the domain logic. It should use the ubiquitous language of the problem domain, and not include any technical implementation of how instances are stored or retrieved from the database, for example.

The api interface also get `client` dependency injections if any side-effects need to be performed, like network calls to services in- or outside the cluster. See the email api's `EmailClient` as an example of how the domain core knows _what_ to do and _when_ to do it, but only the implementation of its interface knows _how_. In this case that is through a postmark client, but this can easily be substituted.

### Config layer
`config` has all the configuration that the microservice persists throughout its lifetime. It is strictly immutable. It includes configuration of the microservice's dependencies to be loosely coupled to its deployment environment.

- Environment variables: the obvious choice to configure a port, graphql endpoint, database configurations, etc.
- Command-line flags: if the microservice accepts command-line flags, they should be loaded in this package, included in the Config variable and returned to the main package.
- Filesystem data: files that are accessible at runtime can be loaded into the service's memory in this package. Mind that it is meant for immutable data. If your files are mutable and service-scoped (probably you don't want this), it should provide a different package to interact with. If it is shared across replicas of this microservice, it's seen as a storage implementation and hence be included therein.

It is currently a _flat_ architecture, in the sense that all packages live in the root, even though some packages are strictly a dependency of others. The intention is to provide a hierarchy as soon as a good solution is found.

The main package that and only that configuration that the other packages need. No other package may load from the inputs that are meant to be in the config layer.

### Server layer

`server` is the source of input for this service for invocations of behaviour. It initialises a server that, in this case, is an http server. Even though the design is meant to be api-agnostic. 

One could, for instance, add a gRPC, REST or SOAP implementation that calls the same domain logic and uses the same storage interface.

Publish/subscribe, event-driven architectures can also be implemented by starting a similar process that takes as input the interface of a domain core; in this case an implementation of the auth.Auth interface. The server in that case is a process that subscribes, listens, polls, or otherwise receives messages and accordingly calls the corresponding methods on the api it has access to.

#### Graph layer

`graph` is an implementation of the graphql resolvers necessary for exposing a graphql api. Only the server calls this package, so it might be included as a `server` subpackage.

### Storage layer

`storage` currently has a postgres database as its storage interface. The repository layer uses this database interface to implement the interface that the domain core needs to perform business logic.

### Repository layer

`repository` sits in between the domain-logic and the storage layer. It translates the retrieval and persistence of entities that need to be committed to storage, and its functionality is called by the domain core package. 

It may either be included as a subpackage under the `server` package or the `auth` domain package, or be their parent. Since it does not include any domain-specific knowledge and mostly implements the interface to query databases, it will probably be included underneath the `server` package.

### Client layer

`client` contains the implementations of any clients that the APIs need. The interfaces are defined there, but its implementations are defined here such that the domain core remains agnostic of any network dependencies and third-party libraries. They are injected as dependencies in the service APIs upon initialisation.

Note that the Repository is effectively a client for a database. For clarity, it is deemed appropriate to have it live on the same level, even though one may claim the client is a generalisation of a repoository.

#### Synchronous calls between APIs

If services do requests or commands to other services, and their dependencies are not circular, one may include a client in one API that makes calls to another. This client can now easily invoke api calls directly, while in the future it can be replaced by a network call if the services live in the same cluster. This is an example of the evolutionary principle discussed above.
