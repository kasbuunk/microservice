package othersvc

// Serves as placeholder for where to put a dependency interface to interact with another service over the network.

// A dependency for another service only requires and returns domain-specific data.

// This interface would not know how the call would be made. Perhaps the request does not even go over the network,
// like the simplified 'email' service is interacted with. At this level, that is unknown.

// These files may be placed in separate packages if it grows, since there should be no dependencies between them.
