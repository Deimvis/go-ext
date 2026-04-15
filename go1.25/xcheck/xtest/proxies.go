package xtest

// TODO: ProxiesToInterface: checks that implementation proxies method calls to interface it implements.
// Usefulness: not rare case when one implements wrap over existing struct and it will be useful to test
// that each method proxies to corresonding method.
// - support ignoring methods (e.g. ProxiesToInterface[io.SeekReader](t, newImpl, IgnoreMethods("Seek")))
// - support specifing args (calls method with empty values by default):
//   ProxiesToInterface[io.SeekReader](t, newImpl, WithMethodArgs("Seek", 0, io.SeekCurrent))

// func ProxiesToInterface[InterfaceT any](t *testing.T, newImpl func(InterfaceT) any, opts...) {}
