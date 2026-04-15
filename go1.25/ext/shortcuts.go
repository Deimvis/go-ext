package ext

// TODO: impl using errgroups ?
// TODO: allow sequential and parallel launches as options?
func UntilFirstErr(fns ...func() error) error {
	for _, f := range fns {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
