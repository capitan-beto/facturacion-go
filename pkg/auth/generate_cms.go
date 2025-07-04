package auth

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"go.mozilla.org/pkcs7"
)

// File reader util.
func readFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// GenerateCMS creates a CMS (Cryptographic Message Syntax) message signed with
// the previous generated PEM Certificate and private key.
func GenerateCMS(TRA []byte, certPath, keyPath string) ([]byte, error) {
	// Load the certificate.
	certPEM, err := readFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("error reading certificate: %v", err)
	}

	// Decode PEM certificate.
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing certificate:", err)
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Load the private key.
	keyPEM, err := readFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %v", err)
	}

	// Decode PEM private key.
	block, _ = pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	// Create PKCS7 signed data.
	signedData, err := pkcs7.NewSignedData(TRA)
	if err != nil {
		return nil, fmt.Errorf("error creating signed data: %v", err)
	}

	// Add signer.
	err = signedData.AddSigner(cert, privateKey, pkcs7.SignerInfoConfig{})
	if err != nil {
		return nil, fmt.Errorf("error adding signer: %v", err)
	}

	// Finalize PKCS7 signature (detached).
	signedBytes, err := signedData.Finish()
	if err != nil {
		return nil, fmt.Errorf("error finishing signed data: %v", err)
	}

	// Encode and return CMS.
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CMS", Bytes: signedBytes})
	return pemData, err
}
