package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	for _, stage := range stages {
		ch = stageWrapper(ch, done, stage)
	}
	return ch
}

func stageWrapper(in In, done In, stage Stage) Out {
	inner := make(Bi)
	ch := stage(inner)

	go func(in In) {
		defer close(inner)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				inner <- v
			case <-done:
				return
			}
		}
	}(in)

	return ch
}
