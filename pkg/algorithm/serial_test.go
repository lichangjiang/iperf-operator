package algorithm

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestSerialize(t *testing.T) {
	nodeHostMap, svcIpMap := createData(5)

	result := Serialize(nodeHostMap, svcIpMap)
	assert.Assert(t, is.Equal(result.JobNodeSize, 20))
	assert.Assert(t, is.Equal(result.EpochSize, 20))

	for i, m := range result.Jobs {
		fmt.Printf("epoch %d -> %+v\n", i, m)
	}
}

func TestParallel(t *testing.T) {
	nodeHostMap, svcIpMap := createData(8)

	parallel := 2

	for parallel <= 8 {
		fmt.Printf("*************parallel:%d****************\n", parallel)
		result := Parallelize(nodeHostMap, svcIpMap, parallel)

		if 56%parallel == 0 {
			assert.Assert(t, is.Equal(result.EpochSize, 56/parallel))
		} else {
			assert.Assert(t, is.Equal(result.EpochSize, 56/parallel+1))
		}

		assert.Assert(t, is.Equal(result.JobNodeSize, 56))
		parallel++
		for i, m := range result.Jobs {
			fmt.Printf("epoch %d -> %+v\n", i, m)
		}
	}

}

func TestFastSerialize(t *testing.T) {

	nodeHostMap, svcIpMap := createData(3)
	result := FastSerialize(nodeHostMap, svcIpMap)
	for i, m := range result.Jobs {
		fmt.Printf("epoch %d -> %+v\n", i, m)
	}
	assert.Assert(t, is.Equal(result.EpochSize, 6))
	assert.Assert(t, is.Equal(result.JobNodeSize, 6))

	fmt.Println("************************************")

	nodeHostMap, svcIpMap = createData(5)
	result = FastSerialize(nodeHostMap, svcIpMap)
	for i, m := range result.Jobs {
		fmt.Printf("epoch %d -> %+v\n", i, m)
	}
	assert.Assert(t, is.Equal(result.EpochSize, 12))
	assert.Assert(t, is.Equal(result.JobNodeSize, 20))

	fmt.Println("************************************")

	nodeHostMap, svcIpMap = createData(8)
	result = FastSerialize(nodeHostMap, svcIpMap)
	for i, m := range result.Jobs {
		fmt.Printf("epoch %d -> %+v\n", i, m)
	}
	assert.Assert(t, is.Equal(result.EpochSize, 14))
	assert.Assert(t, is.Equal(result.JobNodeSize, 56))
}

func createData(size int) (map[string]string, map[string]string) {
	nodeHostMap := map[string]string{}
	svcIpMap := map[string]string{}

	for i := 0; i < size; i++ {
		node := fmt.Sprintf("node%d", i)
		ip := fmt.Sprintf("ip%d", i)
		nodeHostMap[node] = node
		svcIpMap[node] = ip
	}
	return nodeHostMap, svcIpMap
}
