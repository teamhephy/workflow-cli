package parser

import "io"

// HephyCmd is an implementation of Commander.
type FakeHephyCmd struct {
	ConfigFile string
	Warned     bool
	WOut       io.Writer
	WErr       io.Writer
	WIn        io.Reader
}

func (d FakeHephyCmd) Println(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeHephyCmd) Print(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeHephyCmd) Printf(string, ...interface{}) (int, error) {
	return 1, nil
}

func (d FakeHephyCmd) PrintErrln(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeHephyCmd) PrintErr(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeHephyCmd) PrintErrf(string, ...interface{}) (int, error) {
	return 1, nil
}

func (d FakeHephyCmd) ServicesAdd(string, string, string) (error) {
	return nil
}

func (d FakeHephyCmd) ServicesList(string) (error) {
	return nil
}

func (d FakeHephyCmd) ServicesRemove(string, string) (error) {
	return nil
}