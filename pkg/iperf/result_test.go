package iperf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/lichangjiang/iperf-operator/pkg/util"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestHtmlPrint(t *testing.T) {
	content, err := ioutil.ReadFile("../../test/udp_test.json")
	assert.NilError(t, err)

	htmlResult, err := ioutil.ReadFile("../../test/udp_htmlprint_result.html")
	assert.NilError(t, err)

	var iperfJson IperfJson
	assert.NilError(t, json.Unmarshal(content, &iperfJson))

	ics := iperfJson.Analyse()

	server := "192.168.0.1"
	server2 := "10.17.41.11"
	client := "102.168.0.2"
	client1 := "102.168.0.3"
	csKey := CSKey{
		Server: server,
		Client: client,
	}

	csKey1 := CSKey{
		Server: server,
		Client: client1,
	}
	csKey2 := CSKey{
		Server: server2,
		Client: client1,
	}
	csKey3 := CSKey{
		Server: server2,
		Client: client1,
	}
	csKeys := []CSKey{csKey, csKey1}
	csKeys1 := []CSKey{csKey2, csKey3}
	serverKeyMap := map[string][]CSKey{server: csKeys, server2: csKeys1}
	statisMap := map[CSKey]IperfClientStatis{
		csKey:  ics,
		csKey1: ics,
		csKey2: ics,
		csKey3: ics,
	}
	result := HtmlTablePrint(serverKeyMap, statisMap)
	assert.Equal(t, result, string(htmlResult))
	//fmt.Println(result)
}

func TestParseUdp(t *testing.T) {
	content, err := ioutil.ReadFile("../../test/udp_test.json")
	assert.NilError(t, err)
	var iperfJson IperfJson
	assert.NilError(t, json.Unmarshal(content, &iperfJson))
	assert.Equal(t, iperfJson.End.Sum.Packets, 319)
}

func TestJsonEncode(t *testing.T) {
	jsonStr, err := ioutil.ReadFile("../../test/tcp_test.json")
	assert.NilError(t, err)
	b := []byte(jsonStr)
	var iperfJson IperfJson
	assert.NilError(t, json.Unmarshal(b, &iperfJson))
	assert.Assert(t, is.Len(iperfJson.Start.Connected, 1))
	assert.Assert(t, is.Equal(iperfJson.Start.Version, "iperf 3.0.7"))
	assert.Assert(t, is.Len(iperfJson.Intervals, 6))
}

func TestAnalyse(t *testing.T) {
	jsonStr, err := ioutil.ReadFile("../../test/tcp_test.json")
	assert.NilError(t, err)
	iperfJson, err := ParseLog(string(jsonStr))
	assert.NilError(t, err)

	ics := iperfJson.Analyse()
	assert.Assert(t, is.Equal(ics.Version, "iperf 3.0.7"))
	assert.Assert(t, is.Equal(ics.Timestamp, "Tue, 02 Jul 2019 02:36:55 GMT"))
	assert.Assert(t, is.Equal(ics.SystemInfo, "Linux iperf-client-k79rj 4.15.0-50-generic #54-Ubuntu SMP Mon May 6 18:46:08 UTC 2019 x86_64 GNU/Linux\n"))
	assert.Assert(t, is.Equal(ics.IntervalNum, 6))
	assert.Assert(t, is.Equal(len(ics.SocketStatis), 1))

	tt := fmt.Sprintf("%.0f", ics.SocketStatis[0].TotalTrans)
	assert.Assert(t, is.Equal(tt, "103183441898"))
	min := ics.SocketStatis[0].MinBandWidth / 1024 / 1024 / 1024
	max := ics.SocketStatis[0].MaxBandWidth / 1024 / 1024 / 1024
	avg := ics.SocketStatis[0].AvgBandWidth / 1024 / 1024 / 1024
	t.Logf("min:%.2f max:%.2f avg:%.2f", min, max, avg)

	assert.Assert(t, is.Equal(Round2DataStr(ics.SocketStatis[0].MinBandWidth, true), "12.04Gbits/sec"))
	assert.Assert(t, is.Equal(Round2DataStr(ics.SocketStatis[0].MinBandWidth/1024, true), "12.04Mbits/sec"))
	assert.Assert(t, is.Equal(Round2DataStr(ics.SocketStatis[0].MinBandWidth/1024/1024, true), "12.04Kbits/sec"))
	assert.Assert(t, is.Equal(Round2DataStr(ics.SocketStatis[0].MinBandWidth/1024/1024/1024, false), "12.04Bytes"))
}

func TestEmailContent(t *testing.T) {
	jsonStr, err := ioutil.ReadFile("../../test/tcp_test.json")
	assert.NilError(t, err)
	iperfJson, err := ParseLog(string(jsonStr))
	assert.NilError(t, err)

	ics := iperfJson.Analyse()
	key := CSKey{
		Server: "192.168.36.22",
		Client: "192.168.36.22",
	}

	serverKeyMap := map[string][]CSKey{"192.168.36.22": []CSKey{key}}
	statisMap := map[CSKey]IperfClientStatis{key: ics}

	content := HtmlTablePrint(serverKeyMap, statisMap)
	err = util.SendEmail("305120108@qq.com", "iperf result", content)
	assert.NilError(t, err)
}
