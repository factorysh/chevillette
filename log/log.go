package log

type LineReader func(line []byte) (keys []string, err error)
