package iperf

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lichangjiang/iperf-operator/pkg/util"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

const jsonStr = `
{
    "start":    {
        "connected":    [{
                "socket":   4,
                "local_host":   "172.20.14.138",
                "local_port":   39552,
                "remote_host":  "172.20.14.137",
                "remote_port":  9000
            }],
        "version":  "iperf 3.0.7",
        "system_info":  "Linux iperf-client-k79rj 4.15.0-50-generic #54-Ubuntu SMP Mon May 6 18:46:08 UTC 2019 x86_64 GNU/Linux\n",
        "timestamp":    {
            "time": "Tue, 02 Jul 2019 02:36:55 GMT",
            "timesecs": 1562035015
        },
        "connecting_to":    {
            "host": "172.20.14.137",
            "port": 9000
        },
        "cookie":   "iperf-client-k79rj.1562035015.703501",
        "tcp_mss_default":  1398,
        "test_start":   {
            "protocol": "TCP",
            "num_streams":  1,
            "blksize":  131072,
            "omit": 0,
            "duration": 60,
            "bytes":    0,
            "blocks":   0,
            "reverse":  0
        }
    },
    "intervals":    [{
            "streams":  [{
                    "socket":   4,
                    "start":    0,
                    "end":  10.0008,
                    "seconds":  10.0008,
                    "bytes":    17011466218,
                    "bits_per_second":  1.36081e+10,
                    "retransmits":  5,
                    "snd_cwnd": 1054092,
                    "omitted":  false
                }],
            "sum":  {
                "start":    0,
                "end":  10.0008,
                "seconds":  10.0008,
                "bytes":    17011466218,
                "bits_per_second":  1.36081e+10,
                "retransmits":  5,
                "omitted":  false
            }
        }, {
            "streams":  [{
                    "socket":   4,
                    "start":    10.0008,
                    "end":  20.0002,
                    "seconds":  9.9994,
                    "bytes":    18603048960,
                    "bits_per_second":  1.48833e+10,
                    "retransmits":  0,
                    "snd_cwnd": 2373804,
                    "omitted":  false
                }],
            "sum":  {
                "start":    10.0008,
                "end":  20.0002,
                "seconds":  9.9994,
                "bytes":    18603048960,
                "bits_per_second":  1.48833e+10,
                "retransmits":  0,
                "omitted":  false
            }
        }, {
            "streams":  [{
                    "socket":   4,
                    "start":    20.0002,
                    "end":  30.0002,
                    "seconds":  10.0001,
                    "bytes":    17272668160,
                    "bits_per_second":  1.3818e+10,
                    "retransmits":  0,
                    "snd_cwnd": 2801592,
                    "omitted":  false
                }],
            "sum":  {
                "start":    20.0002,
                "end":  30.0002,
                "seconds":  10.0001,
                "bytes":    17272668160,
                "bits_per_second":  1.3818e+10,
                "retransmits":  0,
                "omitted":  false
            }
        }, {
            "streams":  [{
                    "socket":   4,
                    "start":    30.0002,
                    "end":  40.0004,
                    "seconds":  10.0001,
                    "bytes":    16971202560,
                    "bits_per_second":  1.35768e+10,
                    "retransmits":  0,
                    "snd_cwnd": 4203786,
                    "omitted":  false
                }],
            "sum":  {
                "start":    30.0002,
                "end":  40.0004,
                "seconds":  10.0001,
                "bytes":    16971202560,
                "bits_per_second":  1.35768e+10,
                "retransmits":  0,
                "omitted":  false
            }
        }, {
            "streams":  [{
                    "socket":   4,
                    "start":    40.0004,
                    "end":  50.0005,
                    "seconds":  10.0001,
                    "bytes":    17171742720,
                    "bits_per_second":  1.37372e+10,
                    "retransmits":  0,
                    "snd_cwnd": 4203786,
                    "omitted":  false
                }],
            "sum":  {
                "start":    40.0004,
                "end":  50.0005,
                "seconds":  10.0001,
                "bytes":    17171742720,
                "bits_per_second":  1.37372e+10,
                "retransmits":  0,
                "omitted":  false
            }
        }, {
            "streams":  [{
                    "socket":   4,
                    "start":    50.0005,
                    "end":  60.0005,
                    "seconds":  9.99999,
                    "bytes":    16153313280,
                    "bits_per_second":  1.29227e+10,
                    "retransmits":  0,
                    "snd_cwnd": 4203786,
                    "omitted":  false
                }],
            "sum":  {
                "start":    50.0005,
                "end":  60.0005,
                "seconds":  9.99999,
                "bytes":    16153313280,
                "bits_per_second":  1.29227e+10,
                "retransmits":  0,
                "omitted":  false
            }
        }],
    "end":  {
        "streams":  [{
                "sender":   {
                    "socket":   4,
                    "start":    0,
                    "end":  60.0005,
                    "seconds":  60.0005,
                    "bytes":    103183441898,
                    "bits_per_second":  1.37577e+10,
                    "retransmits":  5
                },
                "receiver": {
                    "socket":   4,
                    "start":    0,
                    "end":  60.0005,
                    "seconds":  60.0005,
                    "bytes":    103183441898,
                    "bits_per_second":  1.37577e+10
                }
            }],
        "sum_sent": {
            "start":    0,
            "end":  60.0005,
            "seconds":  60.0005,
            "bytes":    103183441898,
            "bits_per_second":  1.37577e+10,
            "retransmits":  5
        },
        "sum_received": {
            "start":    0,
            "end":  60.0005,
            "seconds":  60.0005,
            "bytes":    103183441898,
            "bits_per_second":  1.37577e+10
        },
        "cpu_utilization_percent":  {
            "host_total":   99.5547,
            "host_user":    1.53012,
            "host_system":  98.0246,
            "remote_total": 6.33031,
            "remote_user":  0.321433,
            "remote_system":    6.00888
        }
    }
}
`

func TestJsonEncode(t *testing.T) {
	b := []byte(jsonStr)
	var iperfJson IperfJson
	assert.NilError(t, json.Unmarshal(b, &iperfJson))
	assert.Assert(t, is.Len(iperfJson.Start.Connected, 1))
	assert.Assert(t, is.Equal(iperfJson.Start.Version, "iperf 3.0.7"))
	assert.Assert(t, is.Len(iperfJson.Intervals, 6))
}

func TestAnalyse(t *testing.T) {
	iperfJson, err := ParseLog(jsonStr)
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
	iperfJson, err := ParseLog(jsonStr)
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
