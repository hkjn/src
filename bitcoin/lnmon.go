// lnmon.go is a simple wrapper around c-lightning's lightning-cli for monitoring state.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type (
	cli          struct{}
	bitcoindInfo struct {
		pid  int
		args []string
	}
	lightningdInfo struct {
		pid  int
		args []string
	}
)

var (
	btcInfo bitcoindInfo
	lnInfo  lightningdInfo
)

func (info bitcoindInfo) String() string {
	if info.pid == 0 {
		return "bitcoindInfo{not running}"
	} else {
		return fmt.Sprintf("bitcoindInfo{pid: %d, args: %q}", info.pid, strings.Join(info.args, " "))
	}
}

func (info lightningdInfo) String() string {
	if info.pid == 0 {
		return "lightningdInfo{not running}"
	} else {
		return fmt.Sprintf("lightningdInfo{pid: %d, args: %q}", info.pid, strings.Join(info.args, " "))
	}
}

// execCmd executes specified command with arguments and returns the output.
func execCmd(cmd string, arg ...string) (string, error) {
	c := exec.Command(cmd, arg...)
	out := bytes.Buffer{}
	stderr := bytes.Buffer{}
	c.Stdout = &out
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			log.Printf("Command %q exited with non-zero status: %v, stderr=%s\n", cmd, err, stderr.String())
			return stderr.String(), nil
		}
		return "", err
	}
	return out.String(), nil
}

func (c cli) exec(cmd string) (string, error) {
	return execCmd("lightning-cli", cmd)
}

func (c cli) GetInfo() (string, error) {
	return c.exec("getinfo")
}

func (c cli) GetPeers() (string, error) {
	return "", nil
}

func (c cli) GetNodes() (string, error) {
	return "", nil
}

func updateBtcInfo() error {
	s, err := execCmd("pgrep", "-a", "bitcoind")
	if err != nil {
		return err
	}
	parts := strings.Split(s, " ")
	if len(parts) > 1 {
		btcInfo.pid = -1
		pid, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		btcInfo.pid = pid

		for _, arg := range parts[1:] {
			btcInfo.args = append(btcInfo.args, arg)
		}
	} else {
		btcInfo = bitcoindInfo{}
	}
	log.Println(btcInfo)
	return nil
}

func updateLnInfo() error {
	s, err := execCmd("pgrep", "-a", "lightningd")
	if err != nil {
		return err
	}
	parts := strings.Split(s, " ")
	if len(parts) > 1 {
		lnInfo.pid = -1
		pid, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		lnInfo.pid = pid

		for _, arg := range parts[1:] {
			lnInfo.args = append(lnInfo.args, arg)
		}
	} else {
		lnInfo = lightningdInfo{}
	}
	log.Println(lnInfo)
	return nil
}

func main() {
	if err := updateBtcInfo(); err != nil {
		log.Fatalf("Failed to get bitcoind info: %v\n", err)
	}
	if err := updateLnInfo(); err != nil {
		log.Fatalf("Failed to get lightningd info: %v\n", err)
	}

	log.Printf("Running lightning-cli getinfo..\n")
	c := cli{}
	info, err := c.GetInfo()
	if err != nil {
		log.Fatalf("Failed to get info: %v\n", err)
	}
	log.Printf("lightningd info: %s\n", info)
	// lightning-cli getinfo
	// lightning-cli getpeers
	// lightning-cli getnodes
	// lightning-cli getchannels

	// lightning-cli listfunds
	// lightning-cli listinvoice
	// lightning-cli listpayments
}
