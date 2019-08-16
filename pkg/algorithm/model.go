package algorithm

type JobNode struct {
	ServerHost string
	ClientHost string
	ServerIp   string
}

type JobMap struct {
	JobNodeSize int
	EpochSize   int
	Jobs        map[int][]JobNode
}
