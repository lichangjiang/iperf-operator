package iperf

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
	KB = 1024

	GbitsSec = "Gbits/sec"
	MbitsSec = "Mbits/sec"
	KbitsSec = "Kbits/sec"
	BitsSec  = "bits/sec"

	GBytes = "GBytes"
	MBytes = "MBytes"
	KBytes = "KBytes"
	Bytes  = "Bytes"

	htmlTableTitle = `
    <tr>
        <td style="text-align:center;font-weight:bold">Server</td>
        <td style="text-align:center;font-weight:bold">Client</td>
        <td style="text-align:center;font-weight:bold">Time</td>
        <td style="text-align:center;font-weight:bold">SocketId</td>
        <td style="text-align:center;font-weight:bold">TotalTransform</td>
        <td style="text-align:center;font-weight:bold">MinBandWidth</td>
        <td style="text-align:center;font-weight:bold">AvgBandWidth</td>
        <td style="text-align:center;font-weight:bold">MaxBandWidth</td>
    </tr>
    `
	htmlTable      = "<table>%s<table>\n"
	htmlCol        = "<td style=\"text-align:center\">%s</td>\n"
	htmlColSpanRow = "<td rowspan=\"%d\">%s</td>\n"
	htmlRow        = "<tr>%s</tr>\n"
)

type Connect struct {
	Socket     int    `json:"socket,omitempty"`
	LocalHost  string `json:"local_host,omitempty"`
	LocalPort  int    `json:"local_port,omitempty"`
	RemoteHost string `json:"remote_host,omitempty"`
	RemotePort int    `json:"remote_port,omitempty"`
}

type Timestamp struct {
	Time     string `json:"time,omitempty"`
	Timesecs int64  `json:"timesecs,omitempty"`
}

type ConnectTo struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type TestStart struct {
	Protocol   string `json:"protocol,omitempty"`
	NumStreams int    `json:"num_streams,omitempty"`
	BlkSize    int64  `json:"blksize,omitempty"`
	Duration   int    `json:"duration,omitempty"`
}

type Stream struct {
	Socket        int     `json:"socket,omitempty"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

type IntervalSum struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

type Sender struct {
	Socket        int     `json:"socket,omitempty"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int     `json:"retransmits"`
}

type Receiver struct {
	Socket        int     `json:"socket,omitempty"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

type EndStream struct {
	Sender   Sender   `json:"sender"`
	Receiver Receiver `json:"receiver"`
}

type SumSent struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int     `json:"retransmits"`
}

type SumReceived struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int     `json:"retransmits"`
}

type CpuUtilPercent struct {
	HostTotal    float64 `json:"host_total"`
	HostUser     float64 `json:"jost_user"`
	HostSystem   float64 `json:"host_system"`
	RemoteTotal  float64 `json:"remote_total"`
	RemoteUser   float64 `json:"remote_user"`
	RemoteSystem float64 `json:"remote_system"`
}

type IperfStart struct {
	Connected     []Connect `json:"connected,omitempty"`
	Version       string    `json:"version"`
	SystemInfo    string    `json:"system_info"`
	Timestamp     Timestamp `json:"timestamp,omitempty"`
	ConnectTo     ConnectTo `json:"connecting_to,omitempty"`
	Cookie        string    `json:"cookie,omitempty"`
	TcpMssDefault int       `json:"tcp_mss_default,omitempty"`
	TestStart     TestStart `json:"test_start,omitempty"`
}

type IperfInterval struct {
	Streams []Stream    `json:"streams"`
	Sum     IntervalSum `json:"sum"`
}

type IperfEnd struct {
	Streams        []EndStream    `json:"streams"`
	SumSent        SumSent        `json:"sum_sent"`
	SumReceived    SumReceived    `json:"sum_received"`
	CpuUtilPercent CpuUtilPercent `json:"cpu_utilization_percent"`
}

type IperfJson struct {
	Start     IperfStart      `json:"start"`
	Intervals []IperfInterval `json:"intervals"`
	End       IperfEnd        `json:"end"`
}

type SocketStatis struct {
	Id           int
	TotalTrans   float64
	MinBandWidth float64
	MaxBandWidth float64
	AvgBandWidth float64
}

type IperfClientStatis struct {
	Start        float64
	End          float64
	IntervalNum  int
	SocketStatis []SocketStatis
	Version      string
	SystemInfo   string
	Timestamp    string
}

type CSKey struct {
	Server string
	Client string
}

func ParseLog(log string) (*IperfJson, error) {
	var j IperfJson
	err := json.Unmarshal([]byte(log), &j)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (j *IperfJson) Analyse() IperfClientStatis {

	version := j.Start.Version
	si := j.Start.SystemInfo
	ts := j.Start.Timestamp.Time
	start := j.End.SumSent.Start
	end := j.End.SumSent.End
	intervalnum := len(j.Intervals)
	sslice := make([]SocketStatis, len(j.End.Streams))
	minMap := make(map[int]float64)
	maxMap := make(map[int]float64)

	for _, interval := range j.Intervals {
		for _, stream := range interval.Streams {
			id := stream.Socket
			min, ok := minMap[id]
			if !ok {
				min = stream.BitsPerSecond
			}

			max, ok := maxMap[id]
			if !ok {
				max = stream.BitsPerSecond
			}

			if stream.BitsPerSecond > max {
				max = stream.BitsPerSecond
			}

			if stream.BitsPerSecond < min {
				min = stream.BitsPerSecond
			}

			minMap[id] = min
			maxMap[id] = max
		}
	}

	for i, stream := range j.End.Streams {
		id := stream.Sender.Socket
		totalTrans := stream.Sender.Bytes
		avgBw := stream.Sender.BitsPerSecond
		min, ok := minMap[id]
		if !ok {
			min = 0.0
		}

		max, ok := maxMap[id]
		if !ok {
			max = 0.0
		}

		ss := SocketStatis{
			Id:           id,
			TotalTrans:   totalTrans,
			AvgBandWidth: avgBw,
			MaxBandWidth: max,
			MinBandWidth: min,
		}

		sslice[i] = ss
	}

	return IperfClientStatis{
		Start:        start,
		End:          end,
		IntervalNum:  intervalnum,
		Version:      version,
		SystemInfo:   si,
		Timestamp:    ts,
		SocketStatis: sslice,
	}
}

func HtmlTablePrint(serverKeyMap map[string][]CSKey,
	statisMap map[CSKey]IperfClientStatis) string {

	var buf bytes.Buffer
	buf.WriteString("<table border=\"1\">\n")

	buf.WriteString(htmlTableTitle)
	for server, csKeys := range serverKeyMap {
		var rowSpan int
		isHead := true

		buf.WriteString("<tr>\n")
		for _, key := range csKeys {
			ics := statisMap[key]
			rowSpan += len(ics.SocketStatis)
		}
		serverCol := fmt.Sprintf(htmlColSpanRow, rowSpan, server)
		buf.WriteString(serverCol)
		for _, key := range csKeys {
			ics := statisMap[key]
			row := ics.htmlRowPrint(key.Client, isHead)
			buf.WriteString(row)
			if isHead == true {
				buf.WriteString("</tr>\n")
				isHead = false
			}
		}
	}

	buf.WriteString("</table>\n")
	return buf.String()
}

func (i IperfClientStatis) htmlRowPrint(client string, isHead bool) string {
	if len(i.SocketStatis) == 0 {
		return ""
	}

	var buf bytes.Buffer
	rowSpan := len(i.SocketStatis)
	if !isHead {
		buf.WriteString("<tr>")
	}
	buf.WriteString(fmt.Sprintf("<td style=\"text-align:center\" rowspan=\"%d\">", rowSpan))
	buf.WriteString(client)
	buf.WriteString("</td>\n")

	buf.WriteString(fmt.Sprintf("<td style=\"text-align:center\" rowspan=\"%d\">", rowSpan))
	buf.WriteString(fmt.Sprintf("%.2f~%.2f<br>IntervalNum:%d", i.Start, i.End, i.IntervalNum))
	buf.WriteString("</td>\n")
	buf.WriteString(i.htmlColPrint())
	if !isHead {
		buf.WriteString("</tr>")
	}
	return buf.String()
}

func (i IperfClientStatis) htmlColPrint() string {
	var buf bytes.Buffer

	for num, ss := range i.SocketStatis {
		id := ss.Id
		totalTran := Round2DataStr(ss.TotalTrans, false)
		min := Round2DataStr(ss.MinBandWidth, true)
		max := Round2DataStr(ss.MaxBandWidth, true)
		avg := Round2DataStr(ss.AvgBandWidth, true)
		if num != 0 {
			buf.WriteString("<tr>")
		}
		buf.WriteString(fmt.Sprintf(htmlCol, fmt.Sprintf("%d", id)))
		buf.WriteString(fmt.Sprintf(htmlCol, totalTran))
		buf.WriteString(fmt.Sprintf(htmlCol, min))
		buf.WriteString(fmt.Sprintf(htmlCol, avg))
		buf.WriteString(fmt.Sprintf(htmlCol, max))
		if num != 0 {
			buf.WriteString("</tr>")
		}
	}
	return buf.String()
}

func Round2DataStr(data float64, isBW bool) string {
	typ, unit := getUnit(data, isBW)
	switch typ {
	case 1:
		return fmt.Sprintf("%.2f %s", data/GB, unit)
	case 2:
		return fmt.Sprintf("%.2f %s", data/MB, unit)
	case 3:
		return fmt.Sprintf("%.2f %s", data/KB, unit)
	default:
		return fmt.Sprintf("%.2f %s", data, unit)
	}
}

func getUnit(data float64, isBW bool) (int, string) {
	if isBW {
		if data > GB {
			return 1, GbitsSec
		} else if data > MB {
			return 2, MbitsSec
		} else if data > KB {
			return 3, KbitsSec
		} else {
			return 4, BitsSec
		}
	} else {
		if data > GB {
			return 1, GBytes
		} else if data > MB {
			return 2, MBytes
		} else if data > KB {
			return 3, KBytes
		} else {
			return 4, Bytes
		}
	}
}
