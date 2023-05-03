package smtp

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func SignDKIM(mail []byte, domain, dkimPrivateKeyPath string) ([]byte, error) {
	r := strings.NewReader(string(mail))
	keyBytes, err := os.ReadFile(dkimPrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dkim priv key")
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed decode PEM")
	}
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PKS8 dkim private key")
	}

	key := parseResult.(*rsa.PrivateKey)
	options := &dkim.SignOptions{
		Domain:   domain,
		Selector: "seekers",
		Signer:   key,
	}

	var b bytes.Buffer
	if err = dkim.Sign(&b, r, options); err != nil {
		return nil, errors.Wrap(err, "failed sign dkim")
	}
	return b.Bytes(), nil
}
