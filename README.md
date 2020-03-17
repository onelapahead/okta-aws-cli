# okta-aws-cli
A CLI tool for authenticating to AWS via Okta, written in Go.

Why another CLI tool for AWS auth via Okta? Because all the others
are written in Java or Python, rather than be made as a standalone
binary which can be easily installed.

Install:

```bash
mkdir -p /tmp/oktaws
curl -sL \
    https://github.com/hfuss/okta-aws-cli/releases/download/v0.0.1-alpha/oktaws_0.0.1-alpha_Darwin_x86_64.tar.gz \
    -o /tmp/oktaws/oktaws.tgz

tar -zxvf /tmp/oktaws/oktaws.tgz -C /tmp/oktaws/
mv -f /tmp/oktaws/oktaws /usr/local/bin/
rm -rf /tmp/oktaws
```

Configure Okta CLI:

```bash

mkdir -p ${HOME}/.okta
cat <<EOF>${HOME}/.okta/config.properties
#OktaAWSCLI
OKTA_ORG=example.okta.com
OKTA_AWS_APP_URL=https://example.okta.com/home/amazon_aws/0oa123fdxc4ghj/272
OKTA_USERNAME=hfuss
OKTA_PASSWORD_CMD=lpass show --password example.okta.com
#OKTA_PASSWORD_CMD=echo "mypassword"
EOF
```

Configure AWS profiles:

```bash
mkdir -p ${HOME}/.aws
cat <<EOF>${HOME}/.aws/config
[profile example-admin]
credential_process = oktaws -account 12345678910 -role AnAdminRole
region = us-east-1
EOF
```

Test the login:

```bash
# currently must be logged into lpass first
lpass login
AWS_PROFILE=example-admin aws sts get-caller-identity
```