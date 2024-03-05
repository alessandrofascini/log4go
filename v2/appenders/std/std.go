package appenders

import "os"

type StdOut struct{}

func (s *StdOut) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func (s *StdOut) Close() error {
	return os.Stdout.Close()
}

type StdErr struct{}

func (s *StdErr) Write(b []byte) (int, error) {
	return os.Stderr.Write(b)
}

func (s *StdErr) Close() error {
	return os.Stderr.Close()
}
