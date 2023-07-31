package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func checkDown(chIn In, done In) Out {
	ch := make(Bi)

	go func() {
		defer close(ch)
		for v := range chIn {
			select {
			case <-done:
				return
			default:
				ch <- v
			}
		}
	}()

	return ch
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(checkDown(in, done))
	}
	return in
}
