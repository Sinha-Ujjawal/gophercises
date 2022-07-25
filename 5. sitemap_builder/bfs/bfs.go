package bfs

type BFSConfig[N comparable] struct {
	neighborFn   func(N) (map[N]bool, error)
	ignoreErrors bool
	maxDepth     uint32
	maxElements  uint64
}

type bfsConfigOpt[N comparable] func(*BFSConfig[N])

func WithIgnoreErrors[N comparable](ignoreErrors bool) bfsConfigOpt[N] {
	return func(bfsConfig *BFSConfig[N]) {
		bfsConfig.ignoreErrors = ignoreErrors
	}
}

func WithMaxDepth[N comparable](maxDepth uint32) bfsConfigOpt[N] {
	return func(bfsConfig *BFSConfig[N]) {
		bfsConfig.maxDepth = maxDepth
	}
}

func WithMaxElements[N comparable](maxElements uint64) bfsConfigOpt[N] {
	return func(bfsConfig *BFSConfig[N]) {
		bfsConfig.maxElements = maxElements
	}
}

func NewBFSConfig[N comparable](
	neighborFn func(N) (map[N]bool, error),
	opts ...bfsConfigOpt[N],
) BFSConfig[N] {
	ret := BFSConfig[N]{
		neighborFn:   neighborFn,
		ignoreErrors: false,
		maxDepth:     uint32(1 << 31),
		maxElements:  uint64(1 << 63),
	}

	for _, opt := range opts {
		opt(&ret)
	}
	return ret
}

type DepthNode[N comparable] struct {
	Node   N
	Depth  uint32
	Parent N
}

func (bfsConfig BFSConfig[N]) BFS(s N) ([]DepthNode[N], error) {
	var ret []DepthNode[N]
	var frontier []DepthNode[N]
	var newFrontier []DepthNode[N]
	seen := map[N]bool{s: true}
	frontier = append(frontier, DepthNode[N]{Node: s, Depth: 0})
	for {
		if frontier == nil {
			break
		}
		ret = append(ret, frontier...)
		if len(ret) == int(bfsConfig.maxElements) {
			break
		}
		newFrontier = nil
		for _, u := range frontier {
			if len(ret)+len(newFrontier) >= int(bfsConfig.maxElements) {
				break
			}
			if u.Depth < bfsConfig.maxDepth {
				neighbors, err := bfsConfig.neighborFn(u.Node)
				if err != nil {
					if !bfsConfig.ignoreErrors {
						return nil, err
					}
				} else {
					for vNode := range neighbors {
						if !seen[vNode] {
							seen[vNode] = true
							v := DepthNode[N]{Node: vNode, Depth: u.Depth + 1, Parent: u.Node}
							newFrontier = append(newFrontier, v)
						}
					}
				}
			}
		}
		frontier = newFrontier
	}
	return ret[:bfsConfig.maxElements], nil
}
