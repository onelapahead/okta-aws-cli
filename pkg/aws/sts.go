package aws

import (
    "github.com/aws/aws-sdk-go/service/sts"
)

func GetToken(roleArn, principalArn, samlAssertion string) string {
    _ = sts.AssumeRoleWithSAMLInput{}
    return ""
}
