package limits_service

// NOTE: nice-to-have, not required.

// keep track of the limits
// when it thinks the user is getting close to a limit, make a request to figure out what the actual
// limits are, and reset to those. If the user is close, issue a warning.
