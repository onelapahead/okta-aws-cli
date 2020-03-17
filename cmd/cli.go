package main

import (
    "github.com/hfuss/okta-aws-cli/v2/pkg/config"
    "log"
)

func main() {

    log.Println(config.Properties.Username)
}
