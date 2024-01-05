package req

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/thank243/iptvChannel/config"
)

func New(conf *config.Config) *Req {
	r := &Req{
		Cli: resty.New().SetRetryCount(3).SetBaseURL(fmt.Sprintf("%s/EPG/jsp", conf.ApiHost)),

		userId:        conf.UserID,
		authenticator: conf.Authenticator,
	}

	return r
}
