package rpc

type VM struct {
}

func (v *VM) Start(_ string, _ *int) error {
	return nil
}

func (v *VM) Stop(_ string, _ *int) error {
	return nil
}

func (v *VM) SetStatus(_ error, _ *int) error {
	return nil
}

func (v *VM) SetReady(_ string, _ *int) error {
	return nil
}

func (v *VM) GetStatus(_ string, _ *error) error {
	return nil
}

func (v *VM) GetCommand(_ string, _ *string) error {
	return nil
}
