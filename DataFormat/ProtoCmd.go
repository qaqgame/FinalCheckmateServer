package DataFormat

const (
	LoginReq uint32 = 1
	LoginRes uint32 = 2
)

const (
	Created = iota
	Started
	Stopped
	Released
)

type Args struct {
	Phase   int
	Phase2  string
}

type Reply struct {
	V       int
}