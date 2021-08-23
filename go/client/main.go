package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/security/advancedtls"
)

const (
	credRefreshInterval = 1 * time.Minute
)

var (
	certFilePath, _ = filepath.Abs("go/testdata/client.pem")
	keyFilePath, _  = filepath.Abs("go/testdata/client-key.pem")
	rootFilePath, _ = filepath.Abs("go/testdata/ca.pem")
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: client HOST:PORT")
	}
	addr := os.Args[1]

	identityOptions := pemfile.Options{
		CertFile:        certFilePath,
		KeyFile:         keyFilePath,
		RefreshDuration: credRefreshInterval,
	}
	identityProvider, err := pemfile.NewProvider(identityOptions)
	if err != nil {
		log.Fatalf("pemfile.NewProvider(%v) failed: %v", identityOptions, err)
	}

	rootOptions := pemfile.Options{
		RootFile:        rootFilePath,
		RefreshDuration: credRefreshInterval,
	}
	rootProvider, err := pemfile.NewProvider(rootOptions)
	if err != nil {
		log.Fatalf("pemfile.NewProvider(%v) failed: %v", rootOptions, err)
	}

	options := &advancedtls.ClientOptions{
		IdentityOptions: advancedtls.IdentityCertificateOptions{
			IdentityProvider: identityProvider,
		},
		VerifyPeer: func(params *advancedtls.VerificationFuncParams) (*advancedtls.VerificationResults, error) {
			return &advancedtls.VerificationResults{}, nil
		},
		RootOptions: advancedtls.RootCertificateOptions{
			RootProvider: rootProvider,
		},
		VType: advancedtls.CertVerification,
	}
	clientTLSCreds, err := advancedtls.NewClientCreds(options)
	if err != nil {
		log.Fatalf("advancedtls.NewClientCreds(%v) failed: %v", options, err)
	}

	kp := keepalive.ClientParameters{
		Time: 1 * time.Minute,
	}

	// Connect in plaintext by specifying grpc.WithInsecure().
	// Don't actually do it because it's a security issue.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(clientTLSCreds), grpc.WithKeepaliveParams(kp))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer conn.Close()

	err = CallBoot(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = CallInfer(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
