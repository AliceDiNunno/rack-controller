package domain

type TracebackEntry struct {
	Filename string
	Line     int
	Method   string
}

type Traceback struct {
	Message   string
	Traceback []TracebackEntry
}
