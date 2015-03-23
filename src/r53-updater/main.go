/*
 * Copyright (C) 2014 David Smith <dizzyd@dizzyd.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package main

import "flag"
import "fmt"
import "os"
import "net/http"
import "io/ioutil"
import "time"

import "github.com/BurntSushi/toml"

import "github.com/mitchellh/goamz/aws"
import "github.com/mitchellh/goamz/route53"

var ARG_CONFIG_FILE string

type Config struct {
	AccessKey string
	SecretKey string
	ZoneID    string
	Name      string
	TTL       int
}


func getIp() (string, error) {
	resp, err := http.Get("http://ipv4.icanhazip.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func usage() {
	fmt.Printf("usage: r53-updater <options>\n")
	fmt.Printf(" options:\n")
	flag.PrintDefaults()
}

func main() {
	// Initialize command-line option parsing
	flag.StringVar(&ARG_CONFIG_FILE, "c", "/etc/r53-updater.config",
		"Configuration file")

	// Process command line flags
	flag.Parse()

	// If something went awry, bail
	if !flag.Parsed() {
		usage()
		os.Exit(-1)
	}

	// Load the config file; bail if that fail
	config := &Config{ TTL: 180 }
	_, err := toml.DecodeFile(ARG_CONFIG_FILE, config)
	if err != nil {
		fmt.Printf("Failed to load config file %s: %+v\n",
			ARG_CONFIG_FILE, err)
		os.Exit(-1)
	}

	// Force TTL to betwen 60 and 3600
	if config.TTL < 60 { config.TTL = 60 }
	if config.TTL > 3600 { config.TTL = 3600}

	// Get the current IP
	ip, err := getIp()
	if err != nil {
		fmt.Printf("Failed to get IP: %+\vn", err)
		os.Exit(-1)
	}

	// Setup AWS auth
	auth := aws.Auth{config.AccessKey, config.SecretKey, ""}
	r53 := route53.New(auth, aws.USEast)

	// Construct a resource record update
	change := &route53.ChangeResourceRecordSetsRequest{
		Comment: "Update",
		Changes: []route53.Change{
			route53.Change{
				Action: "UPSERT",
				Record: route53.ResourceRecordSet{
					Name: config.Name,
					Type: "A",
					TTL: config.TTL,
					Records: []string{ip},
				},
			},
		},
	}

	// Make the change
	_, err = r53.ChangeResourceRecordSets(config.ZoneID, change)
	if err != nil {
		fmt.Printf("Failed to update %s: %+v\n", config.Name, err)
		os.Exit(-1)
	}

	fmt.Printf("%s %s %s", time.Now().Format(time.RFC3339), config.Name, ip);
}
