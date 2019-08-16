package algorithm

func Serialize(nodeHostMap, svcIpMap map[string]string) JobMap {
	return Parallelize(nodeHostMap, svcIpMap, 1)
}

func Parallelize(nodeHostMap, svcIpMap map[string]string, parallel int) JobMap {
	r := []map[string]string{}
	ipHostMap := map[string]string{}
	for _, serverHost := range nodeHostMap {
		serverIp := svcIpMap[serverHost]
		ipHostMap[serverIp] = serverHost
		for _, client := range nodeHostMap {
			if client == serverHost {
				continue
			}

			m := map[string]string{client: serverIp}
			r = append(r, m)
		}
	}

	jobs := map[int][]JobNode{}
	jobMap := JobMap{
		JobNodeSize: 0,
		EpochSize:   0,
		Jobs:        jobs,
	}

	i := 0
	for i < len(r) {
		js := []JobNode{}
		for j := 0; j < parallel; j++ {
			if i+j >= len(r) {
				break
			}
			m := r[i+j]
			for k, v := range m {
				js = append(js, JobNode{
					ServerHost: ipHostMap[v],
					ClientHost: k,
					ServerIp:   v,
				})
			}
		}
		jobMap.Jobs[jobMap.EpochSize] = js
		jobMap.JobNodeSize += len(js)
		jobMap.EpochSize++
		i += parallel
	}
	return jobMap
}
