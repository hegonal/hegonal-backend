package monitorengine

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/shlin168/go-whois/whois"
)

// Start http monitor through goroutine
func startHttpMonitor(httpMonitor models.HttpMonitor) {
	for {
		go checkHttpStatus(httpMonitor)
		time.Sleep(time.Duration(httpMonitor.Interval) * time.Second)
	}
}

// Check http status function
func checkHttpStatus(httpMonitor models.HttpMonitor) {
	req, err := buildRequest(httpMonitor)
	if err != nil {
		log.Error(err)
		return
	}

	client := buildHttpClient(httpMonitor)

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		errorType, err := handleHttpError(err)
		incidentHandler(httpMonitor, errorType, 0, 0, err)
		return
	}
	defer resp.Body.Close()

	ping := time.Since(start).Milliseconds()

	var errorType models.IncidentType
	var expiryDays int

	//check if match config status code
	errorType, err = checkStatusCode(httpMonitor, resp.StatusCode)
	if err != nil {
		incidentHandler(httpMonitor, errorType, 0, resp.StatusCode, err)
		return
	}

	// check if ssl about to expiry
	errorType, expiryDays, err = checkSslExpiry(httpMonitor, resp.TLS)
	if err != nil {
		incidentHandler(httpMonitor, errorType, expiryDays, 0, err)
		return
	}

	// check if domain about to expiry
	errorType, expiryDays, err = checkDomainExpiry(httpMonitor)
	if err != nil {
		incidentHandler(httpMonitor, errorType, expiryDays, 0, err)
		return
	}
	SloveIncidentHandler(httpMonitor, int(ping))
}

// Build request config to request
func buildRequest(httpMonitor models.HttpMonitor) (*http.Request, error) {
	var method string

	switch httpMonitor.HttpMethod {
	case models.HttpGet:
		method = http.MethodGet
	case models.HttpPost:
		method = http.MethodPost
	case models.HttpPut:
		method = http.MethodPut
	case models.HttpPatch:
		method = http.MethodPatch
	case models.HttpDelete:
		method = http.MethodDelete
	case models.HttpHead:
		method = http.MethodHead
	case models.HttpOptions:
		method = http.MethodOptions
	}

	var req *http.Request
	var err error

	if httpMonitor.RequestBody == nil {
		req, err = http.NewRequest(method, httpMonitor.URL, nil)
	} else {
		body := *httpMonitor.RequestBody
		req, err = http.NewRequest(method, httpMonitor.URL, bytes.NewBufferString(body))
	}

	return req, err
}

// Build a custom http(s) client to send requests
func buildHttpClient(httpMonitor models.HttpMonitor) *http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !httpMonitor.CheckSslError},
	}
	client := &http.Client{
		Timeout:   time.Duration(httpMonitor.RequestTimeout) * time.Second,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !httpMonitor.FollowRedirections {
				return http.ErrUseLastResponse
			}
			if httpMonitor.MaxRedirects == 0 {
				return nil
			}
			if len(via) >= httpMonitor.MaxRedirects {
				return fmt.Errorf("stopped after max redirects")
			}
			return nil
		},
	}
	return client
}

// If http error, check what is the error type
func handleHttpError(err error) (errorType models.IncidentType, resultError error) {
	if err, ok := err.(*url.Error); ok && err.Timeout() {
		return models.IncidentTypeTimeout, err
	}

	return models.IncidentTypeOtherError, err
}

// Check status code is match config
func checkStatusCode(httpMonitor models.HttpMonitor, statusCode int) (errorType models.IncidentType, resultError error) {
	statusCategory := getCategory(statusCode)
	if !containsStatusCode(httpMonitor.HttpStatusCodes, statusCategory, statusCode) {
		return models.IncidentTypeStatusCodeNotMatch, fmt.Errorf("status code not match")
	}
	return models.IncidentTypeFine, nil
}

// Check is the ssl license expiry
func checkSslExpiry(httpMonitor models.HttpMonitor, tlsState *tls.ConnectionState) (errorType models.IncidentType, expiryDays int, resultError error) {
	var daysLeft int
	// check if need to check ssl expiry
	if httpMonitor.SslExpiryReminders < 1 {
		return models.IncidentTypeFine, 0, nil
	}

	// check if has ssl
	if tlsState != nil && len(tlsState.PeerCertificates) > 0 {
		certs := tlsState.PeerCertificates
		daysLeft = int(time.Until(certs[0].NotAfter).Hours() / 24)
	} else {
		return models.IncidentTypeFine, 0, nil
	}

	if daysLeft < httpMonitor.SslExpiryReminders {
		log.Info(daysLeft)
		return models.IncidentTypeSSLExpiry, daysLeft, fmt.Errorf("ssl certificate is about to expire")
	}

	return models.IncidentTypeFine, 0, nil
}

// Check is the  domain name about th expiry
func checkDomainExpiry(httpMonitor models.HttpMonitor) (errorType models.IncidentType, daysLeft int, err error) {
	if httpMonitor.DomainExpiryReminders < 1 {
		return models.IncidentTypeFine, 0, nil
	}

	u, err := url.Parse(httpMonitor.URL)
	if err != nil {
		log.Error(err)
		return models.IncidentTypeFine, 0, nil
	}

	domain := u.Hostname()

	ctx := context.Background()
	client, err := whois.NewClient(whois.WithTimeout(time.Duration(httpMonitor.RequestTimeout) * time.Second))
	if err != nil {
		log.Error(err)
		return models.IncidentTypeDomainExpiry, 0, fmt.Errorf("error get this domain's information")
	}

	whoisDomain, err := client.Query(ctx, domain)
	if err != nil {
		log.Error(err)
		return models.IncidentTypeDomainExpiry, 0, fmt.Errorf("error get this domain's information")
	}

	expirationDate, err := time.Parse(time.RFC3339, whoisDomain.ParsedWhois.ExpiredDate)

	if err != nil {
		log.Error(err)
		return models.IncidentTypeDomainExpiry, 0, fmt.Errorf("error get this domain's information")
	}

	daysLeft = int(time.Until(expirationDate).Hours() / 24)
	if daysLeft < httpMonitor.DomainExpiryReminders {
		return models.IncidentTypeDomainExpiry, daysLeft, fmt.Errorf("domain name is about to expire")
	}

	return models.IncidentTypeFine, daysLeft, nil
}

// Convert status number to string
func getCategory(statusCode int) string {
	switch {
	case statusCode >= 100 && statusCode < 200:
		return "1xx"
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500 && statusCode < 600:
		return "5xx"
	default:
		return fmt.Sprintf("%dxx", statusCode/100)
	}
}

// Check if the status code is match respone status code
func containsStatusCode(monitor []string, category string, normaleStatusCode int) bool {
	for _, code := range monitor {
		if strings.EqualFold(code, category) || strings.EqualFold(code, string(normaleStatusCode)) {
			return true
		}
	}
	return false
}