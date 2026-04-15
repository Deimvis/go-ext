package xwrapcall

import (
	"sync/atomic"
)

type StackDebug interface {
	// ActiveStackIndex returns active stack index
	// for exec unit corresponding to given execInd.
	// If no such exec unit or exec unit has InvalidExecInd,
	// it returns false (these cases should not be
	// distinguished by this method).
	ActiveStackIndex(execInd uint64) (StackInd, bool)
}

type stackDebug struct {
	// so far >1 execUnits are not supported,
	// so no synchronization here.
	execs []execUnit
}

type execUnit struct {
	activeStackInd atomic.Int64
}

var _ StackDebug = (*stackDebug)(nil)

func (sb *stackDebug) ActiveStackIndex(execInd uint64) (StackInd, bool) {
	if execInd >= uint64(len(sb.execs)) {
		return InvalidStackInd, false
	}
	ind := sb.execs[execInd].activeStackInd.Load()
	return ind, ind != InvalidStackInd
}
