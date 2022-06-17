# Go Microservice Architecture

This project is licensed under the MIT License.

The microservice in this example exposes a graphql api over http, although the very purpose of this project is to design a microservice to be api-agnostic. 

The domain core is located in the `auth` package, exposing some user-related logic. It explicitly does not import any non-domain package. 

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

This is still a work-in-progress. The packages serve their purpose, but may be renamed and restructured in a suitable hierarchy when one is deemed appropriate.

### What & why

The intent of this project is to provide an example microservice that is scalable in terms of complexity. In order to achieve that, it has implements Domain-Driven Design patterns for separation of concerns.

## Domain-driven-design
The core feature is to showcase a microservice architecture with an inward dependency direction of the following:

### Domain core

In this example, `auth` has any and all domain logic. It should use the ubiquitous language of the problem domain, and not include any technical implementation of how instances are stored or retrieved from the database, for example.

The api interface also has 'broker' clients if any side-effects need to be performed, like network calls to services in- or outside the cluster.

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

