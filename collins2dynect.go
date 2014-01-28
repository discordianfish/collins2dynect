package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/discordianfish/go-collins/collins"
	"github.com/discordianfish/go-dynect/dynect"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	// FIXME: add template for hostname/subdomain
	domain      = flag.String("domain", "example.com", "domain to manage")
	dryRun      = flag.Bool("dry", false, "dry run, do not update dyn")
	dynCustomer = flag.String("dyn.customer", "", "customer")
	dynUser     = flag.String("dyn.user", "", "username")
	dynPass     = flag.String("dyn.pass", "", "password")

	collinsUser = flag.String("collins.user", "blake", "username")
	collinsPass = flag.String("collins.pass", "admin:first", "password")
	collinsUrl  = flag.String("collins.url", "http://localhost:9000/api", "collins api url")

	dynClient *dynect.Client
)

const publishTrue = `{"publish":"true"}`

type dynAllRecords struct {
	Data []string `json:"data"`
}

type dynRecords struct {
	RData map[string]string `json:"rdata"`
}

func deleteAllRecords(domain string, recordType string) error {
	if *dryRun {
		return nil
	}
	resp, err := dynClient.Request("GET", fmt.Sprintf("AllRecord/%s", domain), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	allRecords := &dynAllRecords{}
	if err := json.Unmarshal(body, allRecords); err != nil {
		return err
	}

	for _, path := range allRecords.Data {
		parts := strings.Split(path, "/")
		recordT := parts[2]
		zone := parts[3]
		if recordT != recordType {
			continue
		}
		if zone != domain {
			panic("this shouldn't happen")
		}
		rPath := strings.Join(parts[2:], "/")
		log.Printf("- %s", rPath)
		if err := dynClient.Execute("DELETE", rPath, nil); err != nil {
			return err
		}
	}
	return nil
}

func updateZone(domain string, recordMap map[string][]dynRecords) error {
	for fqdn, records := range recordMap {
		for _, record := range records {
			recordBytes, err := json.Marshal(record)
			if err != nil {
				return err
			}

			path := fmt.Sprintf("ARecord/%s/%s", domain, fqdn)
			log.Printf("+ %s", path)
			if *dryRun {
				continue
			}
			if err := dynClient.Execute("POST", path, bytes.NewReader(recordBytes)); err != nil {
				return err
			}
		}
	}
	log.Printf("Publishing domain %s", domain)
	if *dryRun {
		return nil
	}
	if err := dynClient.Execute("PUT", fmt.Sprintf("Zone/%s", domain), bytes.NewReader([]byte(publishTrue))); err != nil {
		log.Fatal(err)
	}
	return nil
}

func main() {
	flag.Parse()
	if !*dryRun {
		if *dynCustomer == "" || *dynUser == "" || *dynPass == "" {
			log.Printf("Please provide dyn credentials!")
			flag.PrintDefaults()
			os.Exit(1)
		}

		dc, err := dynect.New(*dynCustomer, *dynUser, *dynPass)
		if err != nil {
			log.Fatal(err)
		}
		dynClient = dc
	}
	collinsClient := collins.New(*collinsUser, *collinsPass, *collinsUrl)

	assets, err := collinsClient.FindAllAssets()
	if err != nil {
		log.Fatalf("Couldn't find assets: %s", err)
	}

	if err := deleteAllRecords(*domain, "ARecord"); err != nil {
		log.Fatalf("Couldn't delete old records: %s", err)
	}

	records := map[string][]dynRecords{}
	for _, asset := range assets.Data.Data {
		addresses, err := collinsClient.GetAssetAddresses(asset.Asset.Tag)
		if err != nil {
			log.Fatalf("Couldn't get adresses from collins: %s", err)
		}
		aliases := map[string][]string{}
		for _, alias := range strings.Fields(asset.Attributes["0"]["DNS_ALIASES"]) {
			parts := strings.Split(alias, "@")
			if len(parts) != 2 {
				log.Fatalf("Syntax error when parsing DNS_ALIASES attribute on %s", asset.Asset.Tag)
			}
			name := parts[0]
			pool := parts[1]
			aliases[pool] = append(aliases[pool], name)
		}

		for _, address := range addresses.Data.Addresses {
			pool := strings.ToLower(address.Pool)
			hostname := strings.ToLower(fmt.Sprintf("%s%03d",
				asset.Attributes["0"]["PRIMARY_ROLE"],
				asset.Asset.ID,
			))

			zone := strings.ToLower(strings.Join([]string{
				hostname,
				asset.Attributes["0"]["SECONDARY_ROLE"],
				pool,
				*domain,
			}, "."))

			fqdn := fmt.Sprintf("%s.%s", hostname, zone)
			log.Printf("= %s -> %s:", fqdn, address.Address)
			record := dynRecords{
				RData: map[string]string{
					"address": address.Address,
				},
			}
			records[fqdn] = append(records[fqdn], record)
			if names, ok := aliases[pool]; ok {
				for _, name := range names {
					f := fmt.Sprintf("%s.%s", name, *domain)
					records[f] = append(records[f], record)
				}
			}
		}
	}
	if err := updateZone(*domain, records); err != nil {
		log.Fatalf("Couldn't update zone: %s", err)
	}
}
