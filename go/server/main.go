package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/kou164nkn/grpc-sample/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/security/advancedtls"
)

const (
	credRefreshInterval = 1 * time.Minute
	portNumber          = 13333
)

var (
	certFilePath, _ = filepath.Abs("go/testdata/server.pem")
	keyFilePath, _  = filepath.Abs("go/testdata/server-key.pem")
	rootFilePath, _ = filepath.Abs("go/testdata/ca.pem")
)

func main() {
	identityOptions := pemfile.Options{
		CertFile:        certFilePath,
		KeyFile:         keyFilePath,
		RefreshDuration: credRefreshInterval,
	}
	identityProvider, err := pemfile.NewProvider(identityOptions)
	if err != nil {
		log.Fatalf("pemfile.NewProvider(%v) failed: %v", identityOptions, err)
	}
	defer identityProvider.Close()

	rootOptions := pemfile.Options{
		RootFile:        rootFilePath,
		RefreshDuration: credRefreshInterval,
	}
	rootProvider, err := pemfile.NewProvider(rootOptions)
	if err != nil {
		log.Fatalf("pemfile.NewProvider(%v)failed: %v", rootOptions, err)
	}
	defer rootProvider.Close()

	options := &advancedtls.ServerOptions{
		IdentityOptions: advancedtls.IdentityCertificateOptions{
			IdentityProvider: identityProvider,
		},
		RootOptions: advancedtls.RootCertificateOptions{
			RootProvider: rootProvider,
		},
		RequireClientCert: true,
		VerifyPeer: func(params *advancedtls.VerificationFuncParams) (*advancedtls.VerificationResults, error) {
			fmt.Printf("Client common name: %s.\n", params.Leaf.Subject.CommonName)
			return &advancedtls.VerificationResults{}, nil
		},
		VType: advancedtls.CertVerification,
	}
	serverTLSCreds, err := advancedtls.NewServerCreds(options)
	if err != nil {
		log.Fatalf("advancedtls.NewServerCreds(%v) failed: %v", options, err)
	}

	kep := keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}
	serv := grpc.NewServer(grpc.Creds(serverTLSCreds), grpc.KeepaliveEnforcementPolicy(kep))

	hs := health.NewServer()
	healthgrpc.RegisterHealthServer(serv, hs)
	hs.Resume()

	deepthought.RegisterComputeServer(serv, &Server{})

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println("failed to listen", err)
		os.Exit(1)
	}

	// return after being closed `l`, so there is no need close(l) in the main func.
	serv.Serve(l)
}
