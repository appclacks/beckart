package runner

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
	"github.com/appclacks/beckart/template"
	"github.com/appclacks/beckart/tls"
	"github.com/appclacks/beckart/transformers"
	"go.uber.org/zap"
)

func buildURL(action config.HTTPAction) string {
	return fmt.Sprintf(
		"%s://%s%s",
		action.Protocol,
		net.JoinHostPort(action.Target, fmt.Sprintf("%d", action.Port)),
		action.Path)
}

func isSuccessful(action config.HTTPAction, response *http.Response) bool {
	for _, s := range action.ValidStatus {
		if uint(response.StatusCode) == s {
			return true
		}
	}
	return false
}

func ExecuteHTTP(logger *zap.Logger, store *store.Store, transformers map[string]transformers.Transformer, action config.Action) error {
	tlsConfig, err := tls.GetTLSConfig(action.HTTP.Key, action.HTTP.Cert, action.HTTP.Cacert, action.HTTP.ServerName, action.HTTP.Insecure)
	if err != nil {
		return err
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	ctx := context.Background()
	body, err := template.GenTemplate(store, action.Name, action.HTTP.Body)
	if err != nil {
		return err
	}
	urlTemplate := buildURL(action.HTTP)
	url, err := template.GenTemplate(store, fmt.Sprintf("%s-url", action.Name), urlTemplate)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(action.HTTP.Method, url.String(), body)
	if err != nil {
		return fmt.Errorf("fail to initialize HTTP request: %w", err)
	}
	req.Header.Set("User-Agent", "beckart")
	for k, v := range action.HTTP.Headers {
		header, err := template.GenTemplate(store, fmt.Sprintf("%s-header-%s", action.Name, k), v)
		if err != nil {
			return err
		}
		req.Header.Set(k, header.String())
	}
	if len(action.HTTP.Query) != 0 {
		q := req.URL.Query()
		for k, v := range action.HTTP.Query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	redirect := http.ErrUseLastResponse
	if action.HTTP.Redirect {
		redirect = nil
	}

	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return redirect
		},
	}
	timeout, err := time.ParseDuration(action.HTTP.Timeout)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req = req.WithContext(timeoutCtx)

	for _, transformerRef := range action.Transformers {
		transformer, ok := transformers[transformerRef]
		if !ok {
			return fmt.Errorf("unknown transformer %s", transformerRef)
		}
		logger.Debug(fmt.Sprintf("applying transformer %s", transformerRef))
		err = transformer.Transform(req, store)
		if err != nil {
			return fmt.Errorf("fail to execute transformer %s: %w", transformerRef, err)
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("fail to read request body: %w", err)
	}
	if !isSuccessful(action.HTTP, response) {
		responseBodyStr := string(responseBody)
		maxMessageSize := 1000
		message := responseBodyStr
		if len(responseBodyStr) > maxMessageSize {
			message = responseBodyStr[0:maxMessageSize]
		}
		errorMsg := fmt.Sprintf("HTTP request failed: (status %d) => %s", response.StatusCode, html.EscapeString(message))
		err = errors.New(errorMsg)
		return err
	}
	extractors := action.HTTP.Extractors
	if len(extractors.BodyJSON) > 0 {
		err := extractBody(store, responseBody, extractors.BodyJSON)
		if err != nil {
			return err
		}
	}
	if len(extractors.Headers) > 0 {
		err := extractHeaders(store, response.Header, extractors.Headers)
		if err != nil {
			return err
		}
	}
	if extractors.Body != "" {
		store.Set(extractors.Body, responseBody)
	}
	if extractors.BodyString != "" {
		store.Set(extractors.BodyString, string(responseBody))
	}
	return nil
}
