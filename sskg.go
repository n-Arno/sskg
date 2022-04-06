package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math"
	"os"
	"text/template"
)

func randomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}

var ipsecconf = `# /etc/ipsec.conf - strongSwan IPsec configuration file

# basic configuration

config setup
              strictcrlpolicy=no
              uniqueids = yes
              charondebug = "all"

# Site to Site

conn {{ .Left.Name }}-to-{{ .Right.Name }}
              authby=secret
              left=%defaultroute
              leftid={{ .Left.PublicIP }}
              leftsubnet={{ .Left.PrivateSubnet }}
              leftauth=psk
              right={{ .Right.PublicIP }}
              rightid={{ .Right.PublicIP }}
              rightsubnet={{ .Right.PrivateSubnet }}
              rightauth=psk
              keyexchange=ikev2
              keyingtries=%forever
              fragmentation=yes
              ike=aes192gcm16-aes128gcm16-prfsha256-ecp256-ecp521,aes192-sha256-modp3072
              esp=aes192gcm16-aes128gcm16-ecp256-modp3072,aes192-sha256-ecp256-modp3072
              dpdaction=restart
              auto=route

# /etc/ipsec.secrets - This file holds shared secrets or RSA private keys for authentication.

{{ .Left.PublicIP }} : PSK "{{ .PSK }}"
`

// root structure of the yaml file
type Definition struct {
	Sites []Site `default:"[]" yaml:"sites"`
}

type Site struct {
	Name          string `default:"" yaml:"name"`
	PublicIP      string `default:"" yaml:"public_ip"`
	PrivateSubnet string `default:"" yaml:"private_subnet"`
}

type Data struct {
	Left  Site
	Right Site
	PSK   string
}

// parse default values for sitestructure
func (s *Site) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)

	type plain Site
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <local site> <distant site>\n", os.Args[0])
		os.Exit(1)
	}
	// hardcoded source file name
	siteFile := "sskg.yaml"
	// read content of file
	content, err := ioutil.ReadFile(siteFile)
	if err != nil {
		fmt.Printf("ERROR: %v not readable => %v\n", siteFile, err)
		os.Exit(1)
	}
	// parse yaml
	definition := Definition{}
	err = yaml.Unmarshal([]byte(content), &definition)
	if err != nil {
		fmt.Printf("ERROR: %v not readable => %v\n", siteFile, err)
		os.Exit(1)
	}

	left := Site{}
	found := false
	for i := range definition.Sites {
		if definition.Sites[i].Name == os.Args[1] {
			left = definition.Sites[i]
			found = true
		}
	}
	if !found {
		fmt.Printf("ERROR: %v not found in %v\n", os.Args[1], siteFile)
		os.Exit(1)
	}

	right := Site{}
	found = false
	for i := range definition.Sites {
		if definition.Sites[i].Name == os.Args[2] {
			right = definition.Sites[i]
			found = true
		}
	}
	if !found {
		fmt.Printf("ERROR: %v not found in %v\n", os.Args[2], siteFile)
		os.Exit(1)
	}

	psk := randomBase64String(33)
	tmpl, _ := template.New("ipsec.conf").Parse(ipsecconf)
	data := Data{left, right, psk}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n")
	data = Data{right, left, psk}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
