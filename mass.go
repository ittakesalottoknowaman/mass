package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"mass/props"
	"mass/utils/file"
	"mass/utils/sys"

	"github.com/cihub/seelog"
	uuid "github.com/satori/go.uuid"
)

var fConfig = *flag.String("config", "config.yaml", "config file")

var fIP = *flag.String("ip", "", "ip list file")
var fCommand = *flag.String("cmd", "", "command file")
var fConcurrency = *flag.Int("c", 10, "concurrency number")
var fPort = *flag.String("port", "22", "ssh port")

var fHead = *flag.Int("head", -1, "head")
var fTail = *flag.Int("tail", -1, "tail")

type manager struct {
	id                 string
	config             *props.Config
	command            string
	port               []string
	ipList             []string
	executeSuccessIP   []string
	executeErrorIP     []string
	loginFailIP        []string
	passwordList       []string
	wg                 sync.WaitGroup
	concurrencyChannel chan struct{}
	totalNumber        int32
	executedNumber     int32
}

func new() *manager {
	c, err := props.ParseConfig(fConfig)
	if err != nil {

	}

	return &manager{
		id:                 fmt.Sprintf("%s", uuid.Must(uuid.NewV4())),
		wg:                 sync.WaitGroup{},
		config:             c,
		port:               make([]string, 0, 2),
		ipList:             make([]string, 0, 100),
		executeSuccessIP:   make([]string, 0, 100),
		executeErrorIP:     make([]string, 0, 100),
		loginFailIP:        make([]string, 0, 100),
		concurrencyChannel: make(chan struct{}, fConcurrency),
	}
}

func (m *manager) run() {
	var err error

	// 解析待执行ip列表
	m.ipList, err = parseIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	m.totalNumber = int32(len(m.ipList))

	// 解析执行的命令
	m.command, err = file.ToString(fCommand)
	if err != nil {
		fmt.Println(err)
		return
	}

	m.port = strings.Split(fPort, ",")

	for _, ip := range m.ipList {
		m.concurrencyChannel <- struct{}{}
		m.wg.Add(1)
		atomic.AddInt32(&m.executedNumber, 1)
		percent := int(float32(m.executedNumber) / float32(m.totalNumber) * 100.0)
		fmt.Printf("\r[%s] [%d/%d] %d%%", strings.Repeat("=", percent)+">"+strings.Repeat(" ", 100-percent), m.executedNumber, m.totalNumber, percent)
		go m.execute(ip)
	}

	m.wg.Wait()

	m.writeResult()
}

func (m *manager) writeResult() {
	file.WriteString(filepath.Join(m.config.ResultPath, m.id, "exec_success_ip"), strings.Join(m.executeSuccessIP, "\n"))
	file.WriteString(filepath.Join(m.config.ResultPath, m.id, "exec_error_ip"), strings.Join(m.executeErrorIP, "\n"))
	file.WriteString(filepath.Join(m.config.ResultPath, m.id, "login_fail_ip"), strings.Join(m.loginFailIP, "\n"))
}

func (m *manager) execute(ip string) {
	defer func() {
		<-m.concurrencyChannel
		m.wg.Done()
	}()
	shell := fmt.Sprintf("%s/%s/%s-%d", m.config.ResultPath, m.id, ip, time.Now().UnixNano())
	defer os.Remove(shell)

	port, err := m.scanPort(ip)
	if err != nil {
		m.loginFailIP = append(m.loginFailIP, ip)
		return
	}

	for _, password := range m.passwordList {
		c := fmt.Sprintf("sshpass -p '%s' ssh -tt -p %s sa@%s \"%s\"", password, port, ip, m.command)
		file.WriteString(shell, c)

		str, err := sys.CmdOut("sh", shell)
		if err != nil {
			if err.Error() == "exit status 5" {
				continue
			}
			seelog.Info(fmt.Sprintf("IP:%s\nERROR:%s %s", ip, err, str))
			m.executeErrorIP = append(m.executeErrorIP, ip)
			return
		}

		seelog.Info(fmt.Sprintf("IP:%s\nOUTPUT:%s", ip, str))
		m.executeSuccessIP = append(m.executeSuccessIP, ip)
		return
	}
	m.loginFailIP = append(m.loginFailIP, ip)
}

func (m *manager) scanPort(ip string) (string, error) {
	for _, port := range m.port {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), 3*time.Second)
		if err != nil {
			continue
		}
		return port, nil
	}
	return "", fmt.Errorf("ssh port connection failed")
}

func parseIP() ([]string, error) {
	str, err := file.ToString(fIP)
	if err != nil {
		return nil, err
	}

	// 如果head或tail不为-1，则取对应的ip段
	ipList := strings.Split(str, "\n")

	switch {
	case fHead > len(ipList) || fTail > len(ipList) || fHead > fTail:
		return nil, fmt.Errorf("head or tail input error")
	case fHead != -1 && fTail != -1:
		return ipList[fHead-1 : fTail], nil
	case fHead != -1:
		return ipList[fHead-1:], nil
	case fTail != -1:
		return ipList[:fTail], nil
	default:
		return ipList, nil
	}
}

func main() {
	flag.Parse()
	m := new()
	m.run()
}
