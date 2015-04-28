package mal

import "net/http"

var defaultClient = &http.Client{}

var username, password, userAgent string

func Init(uname, passwd, agent string) {
	username = uname
	password = passwd
	userAgent = agent
}
