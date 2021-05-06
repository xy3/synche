package jobs

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func generateCertAndKey(
	template, parent *x509.Certificate,
	name, crtOut, keyOut string,
	parentPrivateKey *rsa.PrivateKey,
) (privateKey *rsa.PrivateKey, err error) {
	var (
		caBytes       []byte
		rootCaCrtFile afero.File
		rootCaKeyFile afero.File
	)
	if privateKey, err = rsa.GenerateKey(cryptoRand.Reader, 2048); err != nil {
		return
	}
	if parentPrivateKey == nil {
		parentPrivateKey = privateKey
	}
	if caBytes, err = x509.CreateCertificate(cryptoRand.Reader, template, parent, &privateKey.PublicKey, parentPrivateKey); err != nil {
		return
	}
	if rootCaCrtFile, err = files.AppFS.Create(crtOut); err != nil {
		return
	}

	log.Info(color.GreenString("Generated %s certificate: %s", name, crtOut))

	defer rootCaCrtFile.Close()
	if err = pem.Encode(rootCaCrtFile, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
		return
	}

	log.Info(color.GreenString("Generated %s key: %s", name, keyOut))
	if rootCaKeyFile, err = files.AppFS.Create(keyOut); err != nil {
		return
	}
	defer rootCaKeyFile.Close()
	if err = pem.Encode(rootCaKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return
	}
	return
}

func GenerateCertificates(serverCertIPs []string, saveLocation string) (err error) {
	var saveTo = func(name string) string {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(saveLocation, string(os.PathSeparator)), name)
	}

	rootCa := &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63()),
		Subject: pkix.Name{
			Organization: []string{"Synche, INC"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	var rootPrivateKey *rsa.PrivateKey

	// generate root certificates
	if rootPrivateKey, err = generateCertAndKey(rootCa, rootCa, "ca", saveTo("ca.pem"), saveTo("ca.key"), nil); err != nil {
		return
	}

	// generate server certificates
	serverCert := &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63()),
		Subject: pkix.Name{
			Organization: []string{"Synche, INC"},
		},
		IPAddresses: []net.IP{},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}

	for _, ips := range serverCertIPs {
		serverCert.IPAddresses = append(serverCert.IPAddresses, net.ParseIP(ips))
	}

	if _, err = generateCertAndKey(serverCert, rootCa, "server", saveTo("server.pem"), saveTo("server.key"), rootPrivateKey); err != nil {
		return
	}

	// generate client certificates
	clientCert := &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63()),
		Subject: pkix.Name{
			Organization: []string{"Synche, INC"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}

	if _, err = generateCertAndKey(clientCert, rootCa, "client", saveTo("client.pem"), saveTo("client.key"), rootPrivateKey); err != nil {
		return
	}

	return
}
