package aws

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sts"
    "log"
    "time"
)

type Credentials struct {
    sts.Credentials
    Version int `json:"Version"`
}

func GetToken(roleArn, principalArn, samlAssertion string) Credentials {
    mySession := session.Must(session.NewSession())

    svc := sts.New(mySession)
    var expires *int64
    expires = new(int64)
    *expires = time.Hour.Milliseconds() / 1000

    output, err := svc.AssumeRoleWithSAML(&sts.AssumeRoleWithSAMLInput{
        RoleArn: &roleArn,
        PrincipalArn: &principalArn,
        SAMLAssertion: &samlAssertion,
        DurationSeconds: expires,
    })
    if err != nil {
        log.Fatal(err)
    }

    return Credentials{Version: 1, Credentials: *output.Credentials}
}
