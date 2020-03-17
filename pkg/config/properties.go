package config

import (
    "flag"
    "github.com/magiconair/properties"
    "log"
    "os/exec"
    "strings"
)

const configPropertiesFilepath = "${HOME}/.okta/config.properties"

type config struct {
    Organization string `properties:"OKTA_ORG"`
    AwsAppURL    string `properties:"OKTA_AWS_APP_URL"`
    Username     string `properties:"OKTA_USERNAME"`
    AwsRegion    string  `properties:"OKTA_AWS_REGION"`
    PasswordCommand string `properties:"OKTA_PASSWORD_CMD"`
}

var Properties config
var Password   string
var BaseURL    string
var Role       string
var AccountId  string

func init() {
    props := properties.MustLoadFile(configPropertiesFilepath, properties.UTF8)

    if err := props.Decode(&Properties); err != nil {
        log.Fatal(err)
    }

    BaseURL = "https://" + Properties.Organization
    loadPassword()

    role := flag.String("role", "ExampleRole", "The name of the IAM role to assume")
    account := flag.String("account", "123456789", "The ID of the AWS account to use")
    flag.Parse()

    Role = *role
    AccountId = *account
}

func loadPassword() {
    passwordCommand := strings.Split(Properties.PasswordCommand, " ")
    // TODO currently requires `lpass login` to be run first, look into pseudotty
    cmd := exec.Command(passwordCommand[0], passwordCommand[1:]...)
    passwordBytes, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    }
    Password = string(passwordBytes[:len(passwordBytes)-1]) // chomp newline
}
