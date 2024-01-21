package zteg

import (
	"errors"
)

func (c *Client) updateCookie() error {
	resp, err := c.cli.R().SetFormData(map[string]string{
		"UserID":        c.userId,
		"Authenticator": c.authenticator,
	}).Post("platform/auth.jsp")
	if err != nil {
		return err
	}

	cookies := resp.Cookies()
	for i := range cookies {
		if cookies[i].Name == "JSESSIONID" && len(cookies[i].Value) > 0 {
			c.cli.SetCookie(cookies[i])
			return nil
		}
	}

	return errors.New("no valid cookie")
}
