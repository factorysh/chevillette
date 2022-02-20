package input

type LineScanner interface {
	Scan() bool
	Text() string
	Err() error
}
