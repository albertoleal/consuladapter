package consuladapter

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	. "github.com/onsi/gomega"
)

const defaultLogLevel = "info"
const defaultProtocolVersion = 2

const (
	portOffsetDNS = iota
	PortOffsetHTTP
	portOffsetClientRPC
	portOffsetSerfLAN
	portOffsetSerfWAN
	portOffsetServerRPC
	PortOffsetLength
)

type configFile struct {
	BootstrapExpect    int            `json:"bootstrap_expect"`
	Datacenter         string         `json:"datacenter"`
	DataDir            string         `json:"data_dir"`
	LogLevel           string         `json:"log_level"`
	NodeName           string         `json:"node_name"`
	Server             bool           `json:"server"`
	Ports              map[string]int `json:"ports"`
	BindAddr           string         `json:"bind_addr"`
	ProtocolVersion    int            `json:"protocol"`
	StartJoin          []string       `json:"start_join"`
	RetryJoin          []string       `json:"retry_join"`
	RejoinAfterLeave   bool           `json:"rejoin_after_leave"`
	DisableRemoteExec  bool           `json:"disable_remote_exec"`
	DisableUpdateCheck bool           `json:"disable_update_check"`
}

func newConfigFile(
	dataDir string,
	nodeName string,
	clusterStartingPort int,
	index int,
	numNodes int,
) configFile {
	startingPort := clusterStartingPort + PortOffsetLength*index
	ports := map[string]int{
		"dns":      startingPort + portOffsetDNS,
		"http":     startingPort + PortOffsetHTTP,
		"rpc":      startingPort + portOffsetClientRPC,
		"serf_lan": startingPort + portOffsetSerfLAN,
		"serf_wan": startingPort + portOffsetSerfWAN,
		"server":   startingPort + portOffsetServerRPC,
	}

	joinAddresses := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		joinAddresses[i] = fmt.Sprintf("127.0.0.1:%d", clusterStartingPort+i*PortOffsetLength+portOffsetSerfLAN)
	}

	return configFile{
		BootstrapExpect:    numNodes,
		DataDir:            dataDir,
		LogLevel:           defaultLogLevel,
		NodeName:           nodeName,
		Server:             true,
		Ports:              ports,
		BindAddr:           "127.0.0.1",
		ProtocolVersion:    defaultProtocolVersion,
		StartJoin:          joinAddresses,
		RetryJoin:          joinAddresses,
		RejoinAfterLeave:   true,
		DisableRemoteExec:  true,
		DisableUpdateCheck: true,
	}
}

func writeConfigFile(
	configDir string,
	dataDir string,
	nodeName string,
	clusterStartingPort int,
	index int,
	numNodes int,
) string {
	filePath := path.Join(configDir, fmt.Sprintf("%s.json", nodeName))
	file, err := os.Create(filePath)
	Ω(err).ShouldNot(HaveOccurred())

	config := newConfigFile(dataDir, nodeName, clusterStartingPort, index, numNodes)
	configJSON, err := json.Marshal(config)
	Ω(err).ShouldNot(HaveOccurred())

	_, err = file.Write(configJSON)
	Ω(err).ShouldNot(HaveOccurred())

	err = file.Close()
	Ω(err).ShouldNot(HaveOccurred())

	return filePath
}
