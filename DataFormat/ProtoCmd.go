package DataFormat

// LoginReq : info type
const (
	LoginReq uint32 = 1
	LoginRes uint32 = 2
	HeartBeatRequset uint32 = 3
	HeartBeatRsponse uint32 = 4
)

// Args : IPCWord using
type Args struct {
	Phase   int
	Phase2  string
}

// Reply : IPCWork using
type Reply struct {
	V       int
}