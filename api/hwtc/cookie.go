package hwtc

import (
	"errors"
	"regexp"
	"strings"
)

func (c *Client) updateCookie() error {
	// get user token
	token, err := c.getUserToken()
	if err != nil {
		return err
	}
	c.userToken = token

	resp, err := c.cli.R().SetFormData(map[string]string{
		"userToken":     c.userToken,
		"UserID":        c.userId,
		"STBType":       "TY1613",
		"Authenticator": c.authenticator,
	}).Post("ValidAuthenticationHWCTC.jsp")
	if err != nil {
		return err
	}

	// get valid cookie
	cookies := resp.Cookies()
	for i := range cookies {
		if cookies[i].Name == "JSESSIONID" && len(cookies[i].Value) > 0 {
			c.cli.SetCookie(cookies[i])
			return nil
		}
	}

	return errors.New("no valid cookie")
}

func (c *Client) getUserToken() (string, error) {
	resp, err := c.cli.R().SetFormData(map[string]string{
		"UserID": c.userId,
	}).Post("authLoginHWCTC.jsp")
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`userToken.+?"\w+?"`)

	kv := strings.Split(string(re.Find(resp.Body())), "=")
	if len(kv) != 2 {
		return "", errors.New("not found valid user token")
	}

	return strings.Trim(kv[1], "\" "), nil
}
