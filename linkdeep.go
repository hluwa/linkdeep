package main

import (
	"./miner"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strings"
)

var (
	proxyFofa = ""
	proxy     = ""
	scheme    = ""
	host      = ""
	path      = ""
	output    = ""
	command   = &cobra.Command{
		Use:   "linkdeep [target]",
		Short: "Automation discovering from public internet for deeplink.",
		Long: `LinkDeep is a useful tool for discover deeplink from internet.
Enter the uri or facts(scheme, host, path) to generate uri
start deeplink discovering from internet big-data.`,
		RunE: linkdeep,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 && scheme == "" && host == "" && path == "" {
				return errors.New("requires input target uri or set flags for uri fact")
			}
			return nil
		},
	}
)

func initConfig() {
	config := miner.GetConfig()
	if proxy != "" {
		config.Proxy = proxy
	}
	if proxyFofa != "" {
		config.Fofa.Proxy = proxyFofa
	}
}

func generateUri() (uri string) {
	uri = ""
	if scheme != "" {
		uri = fmt.Sprintf("%s://", scheme)
	}
	if host != "" {
		if uri == "" {
			uri = "://"
		}
		uri = fmt.Sprintf("%s%s/", uri, host)
		if path != "" {
			uri = fmt.Sprintf("%s%s", uri, path)
		}
	} else if path != "" {
		uri = fmt.Sprintf("%s/%s", uri, path)
	}
	return
}

func linkdeep(_ *cobra.Command, args []string) error {
	initConfig()

	var uri string
	if len(args) > 0 {
		uri = args[0]
	} else {
		uri = generateUri()
	}

	log.Printf("Start mining for %s\n", uri)
	links, e := miner.FofaMiner(uri)
	if e != nil {
		log.Fatalln(e)
	}

	for _, links := range links {
		log.Println(links)
	}
	if output != "" {
		log.Printf("Save %d link into %s\n", len(links), output)
		return ioutil.WriteFile(output, []byte(strings.Join(links, "\n")), 0777)
	}

	return nil
}

func initCommand() {
	command.Flags().StringVarP(&miner.ConfigPath, "config", "c", "config.json", "config file path")
	command.Flags().StringVarP(&proxyFofa, "proxy-fofa", "", "", "proxy for fofa")
	command.Flags().StringVarP(&proxy, "proxy", "x", "", "proxy for global default")
	command.Flags().StringVarP(&scheme, "scheme", "", "", "scheme for deeplink")
	command.Flags().StringVarP(&host, "host", "", "", "host for deeplink")
	command.Flags().StringVarP(&path, "path", "", "", "path for deeplink")
	command.Flags().StringVarP(&output, "output", "o", "output.txt", "output file path")
}

func main() {
	initCommand()
	_ = command.Execute()
}
