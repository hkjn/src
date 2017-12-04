// client.go implements a GRPC Report client.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	googletime "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "hkjn.me/hkjninfra/telemetry/report"
)

const (
	defaultAddr      = "localhost:50051"
	defaultName      = "world"
	defaultFactsPath = "facts.json"
)

var (
	debugging     = os.Getenv("REPORT_DEBUGGING") == "true"
	addr          = getAddr(defaultAddr)
	tlsCertFile   = os.Getenv("REPORT_TLS_CERT")
	tlsKeyFile    = os.Getenv("REPORT_TLS_KEY")
	tlsCaCertFile = os.Getenv("REPORT_TLS_CA_CERT")
)

func debug(format string, a ...interface{}) {
	if !debugging {
		return
	}
	log.Printf(format, a...)
}

// getAddr returns the address to report in to, given a default.
func getAddr(d string) string {
	addrEnv := os.Getenv("REPORT_ADDR")
	if addrEnv != "" {
		return addrEnv
	}
	return d
}

// getInfo returns the extra info to use when reporting in.
func getInfo(d string) (*pb.ClientInfo, error) {
	factsPath := os.Getenv("REPORT_FACTS_PATH")
	if factsPath == "" {
		factsPath = d
	}
	debug("Reading facts.json from %q..\n", factsPath)
	info := &pb.ClientInfo{}
	f, err := os.Open(factsPath)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(f).Decode(info); err != nil {
		return nil, err
	}
	return info, nil
}

// getClient returns the report client and a func to close the client's connection.
func getClient(addr string) (pb.ReportClient, func() error, error) {
	if tlsCertFile == "" {
		return nil, nil, fmt.Errorf("no TLS cert file set with REPORT_TLS_CERT")
	}
	if tlsKeyFile == "" {
		return nil, nil, fmt.Errorf("no TLS key file set with REPORT_TLS_KEY")
	}
	if tlsCaCertFile == "" {
		return nil, nil, fmt.Errorf("no TLS CA cert file set with REPORT_TLS_CA_CERT")
	}

	// Mutual TLS.
	cert, err := tls.LoadX509KeyPair(
		tlsCertFile,
		tlsKeyFile,
	)

	cp := x509.NewCertPool()
	bs, err := ioutil.ReadFile(tlsCaCertFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ca cert: %s", err)
	}

	ok := cp.AppendCertsFromPEM(bs)
	if !ok {
		return nil, nil, fmt.Errorf("failed to append certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   addr,
		Certificates: []tls.Certificate{cert},
		RootCAs:      cp,
	})

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	return pb.NewReportClient(conn), conn.Close, nil
}

// send reports to the server.
func send(c pb.ReportClient) error {
	info, err := getInfo(defaultFactsPath)
	if err != nil {
		return err
	}
	req := &pb.ReportRequest{
		Ts: &googletime.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   int32(time.Now().Nanosecond()),
		},
		Info: info,
	}
	log.Printf("Sending request: %v\n", req)
	r, err := c.Send(context.Background(), req)
	if err != nil {
		return err
	}
	log.Printf("Got message from server: %q", r.Message)
	ireq := &pb.InfoRequest{}
	iresp, err := c.Info(context.Background(), ireq)
	log.Printf("FIXMEH: Sent inforequest, got reply %v, err %v", iresp, err)
	return nil
}

func main() {
	log.Printf("report_client %s starting..\n", Version)

	log.Printf("Contacting server at tcp addr %q..\n", addr)
	c, close, err := getClient(addr)
	if err != nil {
		log.Fatalf("could not create client: %v\n", err)
	}
	defer close()

	if err := send(c); err != nil {
		log.Fatalf("Could not report: %v", err)
	}
}
