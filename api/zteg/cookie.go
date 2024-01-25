package zteg

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) updateCookie() error {
	resp, err := c.cli.R().SetFormData(map[string]string{
		"UserID":        c.userId,
		"Authenticator": c.authenticator,
	}).Post("platform/auth.jsp")
	if err != nil {
		return err
	}

	// get user token
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return err
	}
	doc.Find("input").Each(func(i int, selection *goquery.Selection) {
		if val, ok := selection.Attr("name"); ok {
			if val == "UserToken" {
				if token, ok := selection.Attr("value"); ok {
					c.userToken = token
				}
			}
		}
	})

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
