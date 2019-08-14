package algorithm

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestSerialize(t *testing.T) {
	nodeHostMap := map[string]string{
		"node1": "node1",
		"node2": "node2",
		"node3": "node3",
		"node4": "node4",
		"node5": "node5",
	}

	svcIpMap := map[string]string{
		"node1": "ip1",
		"node2": "ip2",
		"node3": "ip3",
		"node4": "ip4",
		"node5": "ip5",
	}

	result := Serialize(nodeHostMap, svcIpMap)
	assert.Assert(t, is.Equal(result.JobNodeSize, 20))
	assert.Assert(t, is.Equal(result.EpochSize, 20))

	for i, m := range result.Jobs {
		fmt.Printf("epoch %d -> %+v\n", i, m)
	}
}

func TestParallel(t *testing.T) {
	nodeHostMap := map[string]string{
		"node1": "node1",
		"node2": "node2",
		"node3": "node3",
		"node4": "node4",
		"node5": "node5",
	}

	svcIpMap := map[string]string{
		"node1": "ip1",
		"node2": "ip2",
		"node3": "ip3",
		"node4": "ip4",
		"node5": "ip5",
	}

	parallel := 2

	for parallel <= 21 {
		fmt.Printf("*************parallel:%d****************\n", parallel)
		result := Parallelize(nodeHostMap, svcIpMap, parallel)

		if 20%parallel == 0 {
			assert.Assert(t, is.Equal(result.EpochSize, 20/parallel))
		} else {
			assert.Assert(t, is.Equal(result.EpochSize, 20/parallel+1))
		}

		assert.Assert(t, is.Equal(result.JobNodeSize, 20))
		parallel++
		for i, m := range result.Jobs {
			fmt.Printf("epoch %d -> %+v\n", i, m)
		}
	}

}
