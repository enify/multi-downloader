package worker

type Work interface {
	Do() error
}
