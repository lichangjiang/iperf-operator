package iperf

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	GB = 1000 * 1000 * 1000
	MB = 1000 * 1000
	KB = 1000

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
        <td style="text-align:center;font-weight:bold">TCPMinBandWidth</td>
        <td style="text-align:center;font-weight:bold">TCPAvgBandWidth</td>
        <td style="text-align:center;font-weight:bold">TCPMaxBandWidth</td>
        <td style="text-align:center;font-weight:bold">UDPBandWidth</td>
        <td style="text-align:center;font-weight:bold">UDPJitterDelay</td>
        <td style="text-align:center;font-weight:bold">UDPLostPacket</td>
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
	Packets       int     `json:"packets,omitempty"`
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

type Udp struct {
	Receiver
	JitterMs    float64 `json:"jitter_ms"`
	LostPackets int     `json:"lost_packets"`
	Packets     int     `json:"packets"`
	LostPercent float64 `json:"lost_percent"`
}

type EndStream struct {
	Sender   Sender   `json:"sender,omitempty"`
	Receiver Receiver `json:"receiver,omitempty"`
	Udp      Udp      `json:"udp,omitempty"`
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

type Sum struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         float64 `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	JitterMs      float64 `json:"jitter_ms"`
	LostPackets   int     `json:"lost_packets"`
	Packets       int     `json:"packets"`
	LostPercent   float64 `json:"lost_percent"`
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
	SumSent        SumSent        `json:"sum_sent,omitempty"`
	SumReceived    SumReceived    `json:"sum_received,omitempty"`
	Sum            Sum            `json:"sum,omitempty"`
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

type UdpStatis struct {
	JitterMs      float64
	LostPackets   int
	Packets       int
	LostPercent   float64
	BitsPerSecond float64
}

type IperfClientStatis struct {
	Protocol     string
	Start        float64
	End          float64
	IntervalNum  int
	UdpStatis    UdpStatis
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

func (j *IperfJson) parseUdp() IperfClientStatis {
	version := j.Start.Version
	si := j.Start.SystemInfo
	ts := j.Start.Timestamp.Time
	protocol := j.Start.TestStart.Protocol
	start := j.End.Sum.Start
	end := j.End.Sum.End
	intervalnum := len(j.Intervals)

	return IperfClientStatis{
		Protocol:    protocol,
		Start:       start,
		End:         end,
		IntervalNum: intervalnum,
		Version:     version,
		SystemInfo:  si,
		Timestamp:   ts,
		UdpStatis: UdpStatis{
			BitsPerSecond: j.End.Sum.BitsPerSecond,
			JitterMs:      j.End.Sum.JitterMs,
			LostPackets:   j.End.Sum.LostPackets,
			Packets:       j.End.Sum.Packets,
			LostPercent:   j.End.Sum.LostPercent,
		},
	}
}

func (j *IperfJson) parseTcp() IperfClientStatis {
	version := j.Start.Version
	si := j.Start.SystemInfo
	ts := j.Start.Timestamp.Time
	protocol := j.Start.TestStart.Protocol
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
		Protocol:     protocol,
		Start:        start,
		End:          end,
		IntervalNum:  intervalnum,
		Version:      version,
		SystemInfo:   si,
		Timestamp:    ts,
		SocketStatis: sslice,
	}
}

func (j *IperfJson) Analyse() IperfClientStatis {
	protocol := j.Start.TestStart.Protocol
	if protocol == "UDP" {
		return j.parseUdp()
	} else {
		return j.parseTcp()
	}
}

func HtmlTablePrint(serverKeyMap map[string][]CSKey,
	statisMap map[CSKey]IperfClientStatis) string {

	var buf bytes.Buffer
	buf.WriteString("<table border=\"1\">\n")

	buf.WriteString(htmlTableTitle)
	for server, csKeys := range serverKeyMap {
		rowSpan := len(csKeys)
		isHead := true

		buf.WriteString("<tr>\n")
		serverCol := fmt.Sprintf(htmlColSpanRow, rowSpan, server)
		buf.WriteString(serverCol)
		for _, key := range csKeys {
			if isHead != true {
				buf.WriteString("<tr>\n")
			} else {
				isHead = false
			}
			ics := statisMap[key]
			row := ics.htmlRowPrint(key.Client, isHead)
			buf.WriteString(row)
		}
	}

	buf.WriteString("</table>\n")
	return buf.String()
}

func (i IperfClientStatis) htmlRowPrint(client string, _ bool) string {
	isUdp := false
	if i.SocketStatis == nil || len(i.SocketStatis) == 0 {
		isUdp = true
	}

	var buf bytes.Buffer

	buf.WriteString("<td style=\"text-align:center\">")
	buf.WriteString(client)
	buf.WriteString("</td>\n")

	buf.WriteString("<td style=\"text-align:center\">")
	buf.WriteString(fmt.Sprintf("%.2f~%.2f<br>IntervalNum:%d", i.Start, i.End, i.IntervalNum))
	buf.WriteString("</td>\n")
	buf.WriteString(i.htmlColPrint(isUdp))
	buf.WriteString("</tr>")
	return buf.String()
}

func (i IperfClientStatis) htmlColPrint(isUdp bool) string {
	var buf bytes.Buffer

	if !isUdp {
		ss := i.SocketStatis[0]
		//id := ss.Id
		//totalTran := Round2DataStr(ss.TotalTrans, false)
		min := Round2DataStr(ss.MinBandWidth, true)
		max := Round2DataStr(ss.MaxBandWidth, true)
		avg := Round2DataStr(ss.AvgBandWidth, true)

		//buf.WriteString(fmt.Sprintf(htmlCol, fmt.Sprintf("%d", id)))
		//buf.WriteString(fmt.Sprintf(htmlCol, totalTran))
		buf.WriteString(fmt.Sprintf(htmlCol, min))
		buf.WriteString(fmt.Sprintf(htmlCol, avg))
		buf.WriteString(fmt.Sprintf(htmlCol, max))

		buf.WriteString(fmt.Sprintf(htmlCol, "N/A"))
		buf.WriteString(fmt.Sprintf(htmlCol, "N/A"))
		buf.WriteString(fmt.Sprintf(htmlCol, "N/A"))
	} else {
		buf.WriteString(fmt.Sprintf(htmlCol, "N/A"))
		buf.WriteString(fmt.Sprintf(htmlCol, "N/A"))
		buf.WriteString(fmt.Sprintf(htmlCol, "N/A"))
		//udp数据
		buf.WriteString(fmt.Sprintf(htmlCol, Round2DataStr(i.UdpStatis.BitsPerSecond, true)))
		buf.WriteString(fmt.Sprintf(htmlCol, fmt.Sprintf("%.4f ms", i.UdpStatis.JitterMs)))
		buf.WriteString(fmt.Sprintf(htmlCol,
			fmt.Sprintf("%d/%d (%.2f %%)",
				i.UdpStatis.LostPackets,
				i.UdpStatis.Packets,
				i.UdpStatis.LostPercent)))
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
