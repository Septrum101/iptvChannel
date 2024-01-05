package req

import (
	"github.com/go-resty/resty/v2"
)

type Req struct {
	Cli *resty.Client

	userId        string
	authenticator string
}
