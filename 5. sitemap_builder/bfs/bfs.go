package bfs

import "sync"

type DepthNode[N comparable] struct {
	Node   N
	Depth  uint32
	Parent N
	Err    error
}

type sharedState[N comparable] struct {
	lock        *sync.RWMutex
	processed   map[N]bool
	seen        map[N]bool
	neighborFn  func(N) ([]N, error)
	depth       uint32
	maxDepth    uint32
	newFrontier []DepthNode[N]
	retCh       chan DepthNode[N]
}

func exitWhenAlreadyProcessed[N comparable](u DepthNode[N], state *sharedState[N]) bool {
	state.lock.RLock()
	defer state.lock.RUnlock()
	_, ok := state.processed[u.Node]
	return ok
}

func bfsCallForNode[N comparable](u DepthNode[N], state *sharedState[N]) {
	if exitWhenAlreadyProcessed(u, state) {
		return
	}

	neighbors, err := state.neighborFn(u.Node)

	state.lock.Lock()
	defer state.lock.Unlock()
	state.processed[u.Node] = true
	if err != nil {
		u.Err = err
		state.retCh <- u
		return
	}
	for _, vNode := range neighbors {
		_, ok := state.seen[vNode]
		if !ok {
			state.seen[vNode] = true
			v := DepthNode[N]{Node: vNode, Depth: state.depth + 1, Parent: u.Node}
			state.retCh <- v
			state.newFrontier = append(state.newFrontier, v)
		}
	}
}

func BFS[N comparable](s N, neighborFn func(N) ([]N, error), maxDepth uint32) <-chan DepthNode[N] {
	retCh := make(chan DepthNode[N])
	go func() {
		var wg sync.WaitGroup
		var lock sync.RWMutex
		processed := map[N]bool{}
		seen := map[N]bool{s: true}
		frontier := []DepthNode[N]{{Node: s, Depth: 0}}
		state := sharedState[N]{
			lock:        &lock,
			processed:   processed,
			seen:        seen,
			neighborFn:  neighborFn,
			depth:       0,
			maxDepth:    maxDepth,
			newFrontier: nil,
			retCh:       retCh,
		}
		state.retCh <- frontier[0]
		for depth := uint32(0); depth < maxDepth; depth++ {
			if frontier == nil {
				break
			}
			state.depth = depth
			for _, u := range frontier {
				wg.Add(1)
				go func(u DepthNode[N]) {
					defer wg.Done()
					bfsCallForNode(u, &state)
				}(u)
			}
			wg.Wait()
			frontier = state.newFrontier
			state.newFrontier = nil
		}
		close(retCh)
	}()
	return retCh
}
