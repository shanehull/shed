package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"

	urllib "net/url"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

var url string

var checkcrtCommand = &cli.Command{
	Name:        "checkcrt",
	Description: "A command line tool to retrieve and print a certificate chain info from a specified url.",
	Usage:       "Retrieves and prints a certificate from a specified url.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Aliases:     []string{"u"},
			Value:       "",
			Usage:       "the url to retrieve the cert from",
			Destination: &url,
		},
	},
	Action: func(_ *cli.Context) error {
		var err error

		if url == "" {
			url, err = promptValidUrl()
			if err != nil {
				return err
			}
		}

		printCertChain(getServerCertChain(url))

		return nil
	},
}

func promptValidUrl() (string, error) {
	var val string

	validate := func(input string) error {
		_, err := urllib.Parse(input)
		if err != nil {
			return errors.New("invalid url")
		}

		val = input

		return nil
	}

	s := promptui.Prompt{
		Label:     "Server URL",
		Validate:  validate,
		AllowEdit: true,
	}

	_, err := s.Run()
	if err != nil {
		return "", err
	}

	return val, nil
}

func getServerCertChain(serverName string) [][]*x509.Certificate {
	tlsConfig := &tls.Config{
		ServerName: serverName,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get("https://" + serverName)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return resp.TLS.VerifiedChains
}

func printCertChain(certChains [][]*x509.Certificate) {
	for chainIndex, chain := range certChains {
		fmt.Printf("Chain %d:\n", chainIndex+1)
		for certIndex, cert := range chain {
			fmt.Printf("  Certificate %d:\n", certIndex+1)
			fmt.Printf("    Subject: %s\n", cert.Subject)
			fmt.Printf("    Issuer: %s\n", cert.Issuer)
			fmt.Printf("    Valid From: %s\n", cert.NotBefore.Local())
			fmt.Printf("    Valid Until: %s\n", cert.NotAfter.Local())

			// Print SANs if available
			if len(cert.DNSNames) > 0 {
				fmt.Println("    Subject Alternative Names:")
				for _, name := range cert.DNSNames {
					fmt.Printf("      - %s\n", name)
				}
			}

			fmt.Println("") // Add an empty line for better readability
		}
	}
}
