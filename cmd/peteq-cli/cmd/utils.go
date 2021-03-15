package cmd

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/peteqproj/peteq/pkg/client"
	"gopkg.in/yaml.v2"
)

type (
	clientConfig struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	}
)

func createClientConfiguration() (*client.Configuration, context.Context, error) {
	c := &clientConfig{}
	data, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".peteq/config"))
	if err != nil {
		return nil, nil, err
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, nil, err
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		return nil, nil, err
	}

	cnf := &client.Configuration{
		HTTPClient:    createHTTPClient(),
		DefaultHeader: make(map[string]string),
		UserAgent:     "peteq-cli",
		Debug:         false,
		Scheme:        u.Scheme,
		Servers: client.ServerConfigurations{
			{
				URL: u.Host,
			},
		},
	}
	ctx := context.WithValue(context.Background(), client.ContextAPIKeys, map[string]client.APIKey{
		"ApiKeyAuth": {
			Key: c.Token,
		},
	})
	return cnf, ctx, nil
}

func createHTTPClient() *http.Client {
	customTransport := &(*http.DefaultTransport.(*http.Transport)) // make shallow copy
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient := &http.Client{
		Transport: customTransport,
	}
	return httpClient
}

func storeClientConfiguration(url string, token string) error {
	dir := path.Join(os.Getenv("HOME"), ".peteq")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	c := clientConfig{
		URL:   url,
		Token: token,
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dir, "config"), data, os.ModePerm)
}
