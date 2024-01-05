package req

import (
	"errors"
	"net/http"
)

func (r *Req) updateCookie() error {
	resp, err := r.Cli.R().SetQueryParams(map[string]string{
		"UserID":        r.userId,
		"Authenticator": r.authenticator,
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
		r.Cli.SetCookie(cookie)
		return nil
	}

	return errors.New("no valid cookie")
}
