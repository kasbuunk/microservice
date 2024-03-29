package port

// Serves as placeholder for where to put a localbus interface to interact with another service over the network.

// A localbus for another service only requires and returns domain-specific data.

// This interface would not know how the call would be made. Perhaps the request does not even go over the network,
// like the simplified 'email' service is interacted with. At this level, that is unknown.

// These files may be placed in separate packages if it grows, since there should be no dependencies between them.
