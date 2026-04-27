package gemini

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	//"webd/internal/config"
)

type StatusCode int;
const (
	StatusInput StatusCode = iota 	//10
	StatusSensitiveInput 			//11
	StatusSuccess 					//20
	StatusTempRedirect 				//30
	StatusPermRedirect 				//31
	StatusTempFailure 				//40 45 46 47 48 49
	StatusServerUnavailable 		//41
	StatusCgiError 					//42
	StatusProxyError 				//43
	StatusSlowDown 					//44
	StatusPermFailure				//50 54 55 56 57 58
	StatusNotFound 					//51
	StatusGone 						//52
	StatusProxyRequestRefused 		//53
	StatusBadRequest 				//59
	StatusClientCertificate 		//60 63 64 65 66 67 68 69
	StatusCertificateNotAuthorized	//61
	StatusCertifcateNotValid 		//62
)


type GeminiResponse struct {
	StatusCode int
	Meta string
	Body string
}

func RequestPage(rawUrl string) (string, error) {

	host, path, port, err := parseUrl(rawUrl);
	if err != nil {
		return "", err;
	}

	fmt.Println(net.JoinHostPort(host, port))
	conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), &tls.Config{InsecureSkipVerify: true});
	if err != nil {
		return "", err;
	}

	request := fmt.Sprintf("gemini://%v%v\r\n", host, path);

	_, err = conn.Write([]byte(request));
	if err != nil {
		return "", err;
	}

	certs := conn.ConnectionState().PeerCertificates;
	if len(certs) < 1 {
		return "", fmt.Errorf("no certificate provided by the server");
	}
	cert := certs[0];
	err = verifyFingerprint(host, cert);
	if err != nil {
		return "", err;
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


	return sb.String(), nil;
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

func ParseResponse(response string) (GeminiResponse, error) {

	var r GeminiResponse;

	if len(response) < 1 {
		return r, fmt.Errorf("server returned empty response");
	}

	parts := strings.Split(response, " ");

	statusCode := parts[0];
	converted, err := strconv.Atoi(statusCode);
	if err != nil {
		return r, fmt.Errorf("server returned non-integer status code: %v", err);
	}
	if converted < 10 || converted > 62 {
		return r, fmt.Errorf("server returned a status code that dont exist: %v", converted);
	}
	r.StatusCode = converted;

	var meta string;
	if len(parts) >= 2 {
		meta = strings.TrimSpace(parts[1]);
	}

	if converted == 20 || strings.HasPrefix(statusCode, "2") {
		lines := strings.SplitN(response, "\r\n", 2);
		if len(lines) != 2 {
			return r, fmt.Errorf("server returned success response but without body");
		}
		r.Body = lines[1];
	}

	if strings.HasPrefix(statusCode, "3") {

	}

	if len([]byte(meta)) > 1024 {
		return r, fmt.Errorf("server returned a meta with invalid length: %v", len(meta));
	}
	
	r.Meta = meta;

	return r, nil;
}

