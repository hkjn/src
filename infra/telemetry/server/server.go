// server.go implements a GRPC Report server.
package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	googletime "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb "hkjn.me/hkjninfra/telemetry/report"
)

const defaultAddr = ":50051"

var (
	debugging     = os.Getenv("REPORT_DEBUGGING") == "true"
	slackToken    = os.Getenv("REPORT_SLACK_TOKEN")
	tlsCertFile   = os.Getenv("REPORT_TLS_CERT")
	tlsKeyFile    = os.Getenv("REPORT_TLS_KEY")
	tlsCaCertFile = os.Getenv("REPORT_TLS_CA_CERT")
)

type (
	// clientInfo represents info about a client that has
	// reported to us.
	clientInfo struct {
		// lastSeen is the last time we heard from the client.
		lastSeen time.Time
		// info is extra info reported by the client.
		info *pb.ClientInfo
	}
	// reportServer is used to implement report.GreeterServer.
	reportServer struct {
		// clients is the known clients.
		clients map[string]clientInfo
	}
)

func getAddr(defaultAddr string) string {
	if os.Getenv("REPORT_ADDR") != "" {
		return os.Getenv("REPORT_ADDR")
	}
	return defaultAddr
}

func debug(format string, a ...interface{}) {
	if !debugging {
		return
	}
	log.Printf(format, a...)
}

// newRpcServer returns the GRPC server.
func newRpcServer() (*grpc.Server, error) {
	if tlsCertFile == "" {
		return nil, fmt.Errorf("no TLS cert file set with REPORT_TLS_CERT")
	}
	if tlsKeyFile == "" {
		return nil, fmt.Errorf("no TLS key file set with REPORT_TLS_KEY")
	}
	if tlsCaCertFile == "" {
		return nil, fmt.Errorf("no TLS CA cert file set with REPORT_TLS_CA_CERT")
	}
	cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load x.509 key pair: %v", err)
	}
	cp := x509.NewCertPool()
	bs, err := ioutil.ReadFile(tlsCaCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca cert: %s", err)
	}
	ok := cp.AppendCertsFromPEM(bs)
	if !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	conf := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    cp,
	}

	opts := grpc.Creds(credentials.NewTLS(conf))
	rpcServer := grpc.NewServer(opts)
	s := &reportServer{map[string]clientInfo{}}
	pb.RegisterReportServer(rpcServer, s)
	reflection.Register(rpcServer)
	return rpcServer, nil
}

// sendSlack sends msg to Slack.
func sendSlack(msg string) error {
	slackUrl := "https://hooks.slack.com/services/" + slackToken
	data := struct {
		Text      string `json:"text"`
		LinkNames uint   `json:"link_names"`
		// TODO: Find reason icon_emoji seems to be ignored.
		// IconEmoji string `json:"icon_emoji"`
	}{
		Text:      fmt.Sprintf("`[report_server]` %s", msg),
		LinkNames: 1,
		// IconEmoji: ":heavy_exclamation_mark:",
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode as json: %v\n", err)
		return err
	}
	debug("Sending request to %q: %s\n", slackUrl, b)
	resp, err := http.Post(slackUrl, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Failed to send to Slack: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	debug("Slack replied: %v\n", resp)
	return nil
}

// getTime returns the time.Time equivalent of a timestamp proto message.
func getTime(t *googletime.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func (s *reportServer) Info(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	log.Printf("Received info request\n")
	return &pb.InfoResponse{
		Info: map[string]*pb.ClientInfo{
			"notimplementedyet": &pb.ClientInfo{
				CpuArch:  "gelatinous",
				Hostname: "notimplementedyet-inforesponse",
			},
		},
	}, nil
}

// getInfo describes the client info as a string.
func getInfo(info *pb.ClientInfo) string {
	extra := []string{}
	if info.CpuArch != "" {
		extra = append(extra, fmt.Sprintf("`%s`", info.CpuArch))
	}
	if info.KernelName != "" {
		extra = append(extra, fmt.Sprintf("`%s %s`", info.KernelName, info.KernelVersion))
	}
	if info.Platform != "" {
		extra = append(extra, fmt.Sprintf("`%s`", info.Platform))
	}
	for i := range info.Tags {
		extra = append(extra, fmt.Sprintf("`%s`", info.Tags[i]))
	}
	return fmt.Sprintf("`[%s]` `%s` (%s)", info.Id, info.Hostname, strings.Join(extra, ", "))
}

// Send implements report.ReportServer.
func (s *reportServer) Send(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	c, existed := s.clients[req.Info.Id]
	// TODO: Provide API to fetch clients by hostname
	greeting := "Node"
	if !existed {
		greeting = "New node"
	}
	msg := fmt.Sprintf("%s reported to us: %s", greeting, getInfo(req.Info))
	// TODO: Validate data; seems like it can become corrupt:
	// Full info: allowed_ssh_keys:"memory_total_mb" cpu_arch:"7867"

	log.Println(msg)
	log.Printf("Full info: %+v\n", req.Info)
	if existed {
		log.Printf("Heard from known client for the first time in %v: %s\n", time.Since(c.lastSeen), msg)
	} else {
		log.Printf("Heard from new client: %s\n", msg)
		sendSlack(msg)
	}
	c = clientInfo{
		lastSeen: getTime(req.Ts),
		info:     req.Info,
	}
	s.clients[req.Info.Id] = c
	resp := fmt.Sprintf(
		"Hello %q, thanks for writing me at %v, it is now %v.",
		req.Info.Id,
		getTime(req.Ts),
		time.Now().Unix(),
	)
	log.Printf("Responding to client %q: %q..\n", req.Info.Id, resp)
	return &pb.ReportResponse{Message: resp}, nil
}

func main() {
	addr := getAddr(defaultAddr)
	log.Printf("report_server %s starting, binding at %s..\n", Version, addr)
	if slackToken == "" {
		log.Println("No REPORT_SLACK_TOKEN specified, can't report to Slack.")
	}

	rpcServer, err := newRpcServer()
	if err != nil {
		log.Fatalf("failed to create rpc server: %v\n", err)
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	sendSlack(fmt.Sprintf("%s `report_server` starting..", Version))
	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
