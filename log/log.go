package log

type LineReader func(line []byte) (ip string, err error)
