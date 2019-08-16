package algorithm

import (
	"github.com/RoaringBitmap/roaring"
	"k8s.io/apimachinery/pkg/util/sets"
)

func FastSerialize(nodeHostMap, svcIpMap map[string]string) JobMap {
	jobs := map[int][]JobNode{}
	jobMap := JobMap{
		JobNodeSize: 0,
		EpochSize:   0,
		Jobs:        jobs,
	}

	nodeLen := len(nodeHostMap)
	nodes := make([]string, nodeLen)
	i := 0
	for _, node := range nodeHostMap {
		nodes[i] = node
		i++
	}

	row := 0
	global := roaring.New()
	for jobMap.JobNodeSize < (nodeLen * (nodeLen - 1)) {
		for i := row; i < nodeLen; i++ {
			serverNode := nodes[i]
			serverIp := svcIpMap[serverNode]
			for j := 0; j < nodeLen; j++ {
				clientNode := nodes[j]
				if clientNode == serverNode || global.Contains(uint32(i*nodeLen+j)) {
					continue
				}

				nodeSelected := sets.NewString(serverNode, clientNode)
				js := []JobNode{}

				js = append(js, JobNode{
					ServerHost: serverNode,
					ClientHost: clientNode,
					ServerIp:   serverIp,
				})
				global.Add(uint32(i*nodeLen + j))

				for ii := row + 1; ii < nodeLen; ii++ {
					sn := nodes[ii]
					si := svcIpMap[sn]
					if nodeSelected.Has(sn) {
						continue
					}

					for jj := 0; jj < nodeLen; jj++ {
						cn := nodes[jj]
						if cn == sn || nodeSelected.Has(cn) || global.Contains(uint32(ii*nodeLen+jj)) {
							continue
						}
						nodeSelected.Insert(sn, cn)
						js = append(js, JobNode{
							ServerHost: sn,
							ClientHost: cn,
							ServerIp:   si,
						})

						global.Add(uint32(ii*nodeLen + jj))
						break
					}
				}

				jobMap.Jobs[jobMap.EpochSize] = js
				jobMap.JobNodeSize += len(js)
				jobMap.EpochSize++
			}
		}
		row++
	}

	return jobMap
}
