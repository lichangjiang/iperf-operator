# IperfOperator基于Iperf3的k8s集群自动化网络测试工具
<hr>

## 功能:
- 自动测试k8s集群节点网络联通性
- 自动测试k8s集群节点点对点网络带宽
- 通过邮件发送结果报告

## 使用：
1. 部署IperfTask CRD资源，IperfTask用于描述一次测试任务:<br>
`kubectl apply -f ./deploy/crd/ipefTask.yaml`

2. 创建命名空间与邮件服务器信息secret: <br>
>注意k8s secret字符串要求base64编码

```
apiVersion: v1
kind: Namespace
metadata:
  name: iperf-operator
---
apiVersion: v1
kind: Secret
metadata:
  namespace: iperf-operator
  name: iperf-email-secret
type: Opaque
data:
  user: YOUR_EMAIL_USERNAME
  password: YOUR_EMAIL_PWD
  smtp: YOUR_EMAIL_SMTP_SERVER_ADDRESS
  port: SMTP_SERVER_PORT
```
3. 部署IperfOperator，IperfOperator负责监控IperfTask并执行任务:<br>
`kubectl apply -f ./deploy/iperf_operator.yaml`

4. 部署一次IperfTask测试任务(默认执行tcp测试带宽):
```
apiVersion: iperf.test.svc/alpha1
kind: IperfTask
metadata: 
  name: my-iperf
  namespace: iperf-operator
spec:
  iperfImage: networkstatic/iperf3
  toEmail: YOUR_RECEIVE_EMAIL
  serverSpec:
    port: 9000
  clientSpec:
    mode: "fast"
    interval: 2
    duration: 10
```
- interval:等于Iperf3的-i参数，测试间隔
- duration:等于iperf3的测试持续时间
- mode:分为fast和low两种模式，low模式是节点完全点对点测试，fast模式是节点快速点对点测试。

5. 要部署udp测试延迟和丢包率任务执行:
```
apiVersion: iperf.test.svc/alpha1
kind: IperfTask
metadata: 
  name: my-udp-iperf
  namespace: iperf-operator
spec:
  iperfImage: networkstatic/iperf3
  toEmail: YOUR_RECEIVE_EMAIL
  serverSpec:
    port: 9000
  clientSpec:
    udp: true
    bwLimit: "1000M"
    mode: "fast"
    interval: 10
    duration: 60
```
- bwLimit: 限制udp带宽参数，等于iperf3的-b参数。

## 结果报告:
**tcp带宽测试报告：**
![tcp结果邮件报告](https://github.com/lichangjiang/iperf-operator/blob/master/image/email_report.png)

**udp延迟与丢包率测试报告**
![udp结果邮件报告](https://github.com/lichangjiang/iperf-operator/blob/master/image/udp_email_report.png)

## Licensing

Iperf-operator is under the Apache 2.0 license.
