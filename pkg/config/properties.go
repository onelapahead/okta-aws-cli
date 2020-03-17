package config

import (
    "github.com/magiconair/properties"
    "log"
    "os/exec"
    "strings"
)

//#OktaAWSCLI
//OKTA_ORG=bandwidth.okta.com
//OKTA_AWS_APP_URL=https://bandwidth.okta.com/home/amazon_aws/0oa1f31om0gDQiIkh1d8/272
//OKTA_USERNAME=hfuss
//OKTA_BROWSER_AUTH=false
//OKTA_ENV_MODE=true
//#OKTA_MFA_CHOICE=GOOGLE.token:software:totp
//OKTA_STS_DURATION=43200
//OKTA_AWS_REGION=us-east-1
//OKTA_AWS_ROLE_TO_ASSUME=arn:aws:iam::103854071333:role/BWApp_keystone_DevAdminAccess
//OKTA_PASSWORD_CMD=lpass show --password bandwidth.okta.com
//#OKTA_PASSWORD_CMD=echo "mypassword"


const configPropertiesFilepath = "${HOME}/.okta/config.properties"

type config struct {
    Organization string `properties:"OKTA_ORG"`
    AwsAppURL    string `properties:"OKTA_AWS_APP_URL"`
    Username     string `properties:"OKTA_USERNAME"`
    //BrowserAuth  bool   `properties:"OKTA_BROWSER_AUTH"`
    //EnvMode      bool   `properties:"OKTA_ENV_MODE"`
    //MfaChoice    string `properties:"OKTA_MFA_CHOICE"`
    //StsDuration  time.Duration `properties:"OKTA_STS_DURATION,default=3600s"`
    AwsRegion    string  `properties:"OKTA_AWS_REGION"`
    //AwsRoleToAssume string `properties:"OKTA_AWS_ROLE_TO_ASSUME"`
    PasswordCommand string `properties:"OKTA_PASSWORD_CMD"`
}

var Properties config
var Password   string
var BaseURL    string

func init() {
    props := properties.MustLoadFile(configPropertiesFilepath, properties.UTF8)

    if err := props.Decode(&Properties); err != nil {
        log.Fatal(err)
    }

    BaseURL = "https://" + Properties.Organization
    loadPassword()
}

func loadPassword() {
    passwordCommand := strings.Split(Properties.PasswordCommand, " ")
    // currently requires `lpass login` to be run first, look into pseudotty
    cmd := exec.Command(passwordCommand[0], passwordCommand[1:]...)
    passwordBytes, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    }
    Password = string(passwordBytes[:len(passwordBytes)-1]) // chomp newline
}
