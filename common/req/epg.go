package req

import (
	"strconv"
	"strings"
	"time"
)

func (r *Req) GetEPGBytes(channelId int) ([]byte, error) {
	var buf []byte
	for i := 0; i < 3; i++ {
		resp, err := r.Cli.R().ForceContentType("text/html;charset=UTF-8").SetQueryParam("channelId", strconv.Itoa(channelId)).
			Get("stliveplay_30/en/getTvodData.jsp")
		if err != nil {
			return nil, err
		}

		if strings.Contains(resp.String(), "(\"resignon\",\"1\")") {
			time.Sleep(time.Second * 3)
			if err := r.updateCookie(); err != nil {
				return nil, err
			}

			continue
		}
		buf = resp.Body()
		break
	}
	return buf, nil
}
