package restserver

// Serves as placeholder for where to put an input source that is available over the gRPC protocol.
// Please refer to gqlserver for an example of how this would look.
//
// The key idea here is to show that such a server implementation only needs to worry about the
// transport of requests and their types. It translates requests and does little else than invoking
// the behaviour defined in the API interfaces.
//
// One other role may be to choose what errors to show and filter, such that the developer can choose full
// transparency towards services that are internal to the application and only show 'user errors' to the
// caller if e.g. only one particular api is exposed to external consumers.
