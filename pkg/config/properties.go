package config

import (
    "github.com/magiconair/properties"
    "log"
    "os"
)

const configProperties = "${HOME}/.okta/config.properties"

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

func init() {
    props := properties.MustLoadFile(configProperties, properties.UTF8)

    if err := props.Decode(&Properties); err != nil {
        log.Println("Failed to deserialize config {}", err)
        os.Exit(1)
    }
}
