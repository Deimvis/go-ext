package xtesting

// TODO: maybe move to xsync as DisablableBarrier (/SwitchableBarrier/OptionalBarrier)

// --- LEGACY ---

// TODO: maybe move to xsync as Barrier and RecordingGate
// - need to make sure tags approach the best abstraction for identifying parkers
// - not sure about "Distinguishable" or "HeterogeneoousParkers" in naming to emphasize that parkers are identifiable by tags

// # Barrier vs Gate
// - Barrier
//   - specifics
//     - eventually no one reaches it and it is closed
//     - parkers are inspected on park
//   - usage examples:
//     - check what units parked
//     - pass only specific units
//   - closest equivalent:
//     - bidirectional channel
// - Gate
//   - specifics
//     - eventually it is open and anyone passes through it
//     - parkers are inspected after passing
//   - usage examples:
//     - keep units in the part before gate
//     - wait on units passed
//   - closest equivalent:
//     - no-op closable channel + wait group
// - rationale for having them:
//   - have more specific and limited interface for common synchronization techniques used in testing making code more readable and less vulnerable to making mistakes
