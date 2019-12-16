package taskpool

var (
	SignalShutdown = Control{0}
)

type Control struct {
	OpCode int
}
