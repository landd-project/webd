package gemini

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"webd/internal/config"
)

func RequestPage(rawUrl string) error {

	host, path, port, err := parseUrl(rawUrl);
	if err != nil {
		return err;
	}

	fmt.Println(net.JoinHostPort(host, port))
	conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), &tls.Config{InsecureSkipVerify: true});
	if err != nil {
		return err;
	}

	request := fmt.Sprintf("gemini://%v%v\r\n", host, path);

	_, err = conn.Write([]byte(request));
	if err != nil {
		return err;
	}

	certs := conn.ConnectionState().PeerCertificates;
	if len(certs) < 1 {
		return fmt.Errorf("no certificate provided by the server");
	}
	cert := certs[0];
	err = verifyFingerprint(host, cert);
	if err != nil {
		return err;
	}

	buf := make([]byte, 4096);

	var sb strings.Builder;

	for {
		n, err := conn.Read(buf);
		if n > 0 {
			sb.Write(buf[:n]);
		}
		if err != nil {
			break;
		}
	}

	fmt.Println(sb.String());

	return nil;
}

func parseUrl(rawUrl string) (string, string, string, error) {

	if !strings.HasPrefix(rawUrl, "gemini://") {
		rawUrl = fmt.Sprintf("gemini://%v", rawUrl);
	}
	
	u, err := url.Parse(rawUrl);
	if err != nil {
		return "", "", "", err;
	}

	path := u.Path;
	if path == "" {
		path = "/"
	}

	port := u.Port();
	if port == "" {
		port = "1965"
	}

	return u.Hostname(), path, port, nil;
}

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

