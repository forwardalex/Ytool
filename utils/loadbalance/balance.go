package loadbalance

type Backend struct {
	Name    string
	Weight  int
	Current int
	Info    interface{}
}

type BackendList []Backend

// SmoothBalance 平滑负载均衡 参考bfe  current 初始化为权重
func SmoothBalance(backs BackendList) (*Backend, error) {
	var best *Backend
	total, max := 0, 0

	for _, backend := range backs {

		// select backend with greatest current weight
		if best == nil || backend.Current > max {
			best = &backend
			max = backend.Current
		}
		total += backend.Current

		// update current weight
		backend.Current += backend.Weight
	}

	// update current weight for chosen backend
	best.Current -= total

	return best, nil
}
