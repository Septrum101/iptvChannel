package hwtc

import (
	"errors"
	"net/http"
)

func (c *Client) updateCookie() error {
	resp, err := c.cli.R().SetQueryParams(map[string]string{
		"UserID":        c.userId,
		"Authenticator": c.authenticator,
	}).Post("ValidAuthenticationHWCTC.jsp")
	if err != nil {
		return err
	}

	var isLogin bool
	cookie := new(http.Cookie)

	cookies := resp.Cookies()
	for i := range cookies {
		if cookies[i].Name == "JSESSIONID" && len(cookies[i].Value) > 0 {
			cookie = cookies[i]
			isLogin = true
		}
	}
	if isLogin {
		c.cli.SetCookie(cookie)
		return nil
	}

	return errors.New("no valid cookie")
}
