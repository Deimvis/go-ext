package xsync

// TODO: freezable over abitrary struct with IConst + Imutable interfaces
//       wrapper will have methods: Const(func(IConst)) + Mutable(func(IMutable))
//       will work like this:
//         state := xsync.AsFreezable[ConstState, MutableState](&State{})
//         value := state.Const(ConstState.GetValue)
//         state.Mutable(MutableState.IncValue)
//         var old int
//         state.Mutable(func(s MutableState) {
// 				old = s.SetValue(42)
//         })

type Freezable[T any] interface {
	// Freeze returns false when it is already freezed.
	Freeze(violationCallback func()) bool
	Freezed() bool
}

// TODO: Unfreezable ? requires KeepFreezed(func()) method for synchronization
