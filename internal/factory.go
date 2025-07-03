package internal

type TempooFactory struct {
	instance *Tempoo
}

func NewTempooFactory() (*TempooFactory, error) {
	instance, err := NewTempoo()
	if err != nil {
		return nil, err
	}
	return &TempooFactory{instance: instance}, nil
}

func (f *TempooFactory) GetClient() *Tempoo {
	return f.instance
}
