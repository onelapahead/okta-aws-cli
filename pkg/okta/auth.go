package okta

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/RobotsAndPencils/go-saml"
    "github.com/hfuss/okta-aws-cli/v2/pkg/aws"
    "github.com/hfuss/okta-aws-cli/v2/pkg/config"
    "golang.org/x/net/html"
    "golang.org/x/net/publicsuffix"
    "log"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"
    "time"
)

type userCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type sessionRequest struct {
    SessionToken string `json:"sessionToken"`
}

var client *http.Client

func init() {
    jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
    if err != nil {
        log.Fatal(err)
    }

    client = &http.Client{
        Jar: jar,
    }
}

func getSessionToken() string {
    creds := userCredentials{
        Username: config.Properties.Username,
        Password: config.Password,
    }

    credsJson, err := json.Marshal(creds)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Post(config.BaseURL + "/api/v1/authn", "application/json", bytes.NewBuffer(credsJson))
    if err != nil {
        log.Fatal(err)
    }

    if resp.StatusCode != 200 {
        log.Fatal(errors.New("invalid response code " + resp.Status))
    }

    var authResp interface{}

    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        log.Fatal(err)
    }

    auth := authResp.(map[string]interface{})

    switch auth["status"].(string) {
    case "SUCCESS":
        return auth["sessionToken"].(string)
    }
    log.Fatal("Invalid status " + auth["status"].(string))
    return ""
}


func getSessionId(sessionToken string) string {
    sessionReq := sessionRequest{SessionToken:sessionToken}

    sessionJson, err := json.Marshal(sessionReq)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Post(config.BaseURL + "/api/v1/sessions", "application/json", bytes.NewBuffer(sessionJson))
    if err != nil {
        log.Fatal(err)
    }

    if resp.StatusCode != 200 {
        log.Fatal(errors.New("invalid response code " + resp.Status))
    }

    var sessionResp interface{}

    if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
        log.Fatal(err)
    }

    session := sessionResp.(map[string]interface{})
    return session["id"].(string)
}

func addCookie(cookieName, cookie string) {
    baseUrl, err := url.Parse(config.BaseURL)
    if err != nil {
        log.Fatal(err)
    }

    cookies := client.Jar.Cookies(baseUrl)

    sidCookie := &http.Cookie{Value: cookie, Name: cookieName, Expires: time.Now().Add(10 * time.Hour)}
    newCookies := append(cookies, sidCookie)

    client.Jar.SetCookies(baseUrl, newCookies)
}

func getAttrMap(node *html.Node) map[string]string {
    attrs := make(map[string]string, len(node.Attr))
    for _, attr := range node.Attr {
        attrs[attr.Key] = attr.Val
    }
    return attrs
}

func getSamlAssertion() string {
    resp, err := client.Get(config.Properties.AwsAppURL)
    if err != nil {
        log.Fatal(err)
    }

    if resp.StatusCode != 200 {
        log.Fatal(errors.New("invalid response code " + resp.Status))
    }

    doc, err := html.Parse(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    var samlAssertion string
    var samlCrawler func(*html.Node)
    samlCrawler = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "input" {
            attrMap := getAttrMap(n)
            if attrMap["name"] == "SAMLResponse" {
                samlAssertion = attrMap["value"]
                return
            }
        }
        for child := n.FirstChild; child != nil; child = child.NextSibling {
            samlCrawler(child)
        }
    }
    samlCrawler(doc)

    return samlAssertion
}

func extractPrincipalAndRoleArns(samlAssertion string) (string, string) {
    // TODO consider other SAML parsers
    response, err := saml.ParseEncodedResponse(samlAssertion)
    if err != nil {
        log.Fatal(err)
    }

    for _, attr := range response.Assertion.AttributeStatement.Attributes {
        if attr.Name == "https://aws.amazon.com/SAML/Attributes/Role" {
            for _, val := range attr.AttributeValues {
                // TODO better parsing and validation of role
               if strings.Contains(val.Value, config.AccountId) && strings.Contains(val.Value, config.Role) {
                   arns := strings.Split(val.Value, ",")
                   return arns[0], arns[1]
               }
            }
        }
    }
    log.Fatal("Could not match role")
    return "", ""
}

func LoginAws() {
    sessionToken := getSessionToken()
    sessionId := getSessionId(sessionToken)
    addCookie("sid", sessionId)
    assertion := getSamlAssertion()
    principalArn, roleArn := extractPrincipalAndRoleArns(assertion)

    creds := aws.GetToken(roleArn, principalArn, assertion)
    credsJson, err := json.Marshal(&creds)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(credsJson))
}
