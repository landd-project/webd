package gemini

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"webd/internal/config"
)

func saveFingerprint(host string, fingerprint string) error {
	dir, err := config.GetBaseDir();
	if err != nil {
		return err;
	}

	path := filepath.Join(dir, "known");

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644);
	if err != nil {
		return err;
	}
	defer f.Close();

	_, err = fmt.Fprintf(f, "%v %v\n", host, fingerprint);
	if err != nil {
		return err;
	}

	return nil;
}

func verifyFingerprint(host string, cert *x509.Certificate) error {
	sum := sha256.Sum256(cert.Raw);
	fingerprint := hex.EncodeToString(sum[:]);

	dir, err := config.GetBaseDir();
	if err != nil {
		return err;
	}

	path := filepath.Join(dir, "known");

	known := make(map[string]string);

	content, err := config.ReadFileToString(path);
	if err != nil {
		return err;	
	}

	lines := strings.SplitSeq(content, "\n");
	for l := range lines {
		fields := strings.Fields(l);
		if len(fields) == 2{
			known[fields[0]] = fields[1];
		}
	}

	if saved, exists := known[host]; exists {
		if saved != fingerprint {
			return fmt.Errorf("fingerprint for `%v` has changed", host);
		}
		return nil;
	}

	err = saveFingerprint(host, fingerprint);
	if err != nil {
		return err;
	}

	return nil;
}

