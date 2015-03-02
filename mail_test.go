package go_mailgun_client

import (
	"testing"
)

const (
	apiKey  = "<Input your API Key>"
	domain  = "<Input your domain in mailgun>"
	from    = "<Input from address>"
	to      = "<Input to address>"
	subject = "メールガンから送信"
	body    = "テストコードから本文"
)

func Test_OK(t *testing.T) {
	client := New(apiKey, domain, from)
	err := client.Send(to, subject, body)
	if err != nil {
		t.Errorf("Unexepcted Error : %s", err)
	}
}
