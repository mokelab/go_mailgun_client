package go_mailgun_client

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

type client struct {
	apiKey string
	domain string
	from   string
}

type MailgunError struct {
	Status int
	Body   string
}

func New(apiKey, domain, from string) *client {
	return &client{
		apiKey: apiKey,
		domain: domain,
		from:   from,
	}
}

func (c *client) Send(email, subject, body string) error {
	url := "https://api.mailgun.net/v2/" + c.domain + "/messages"
	// build parameters (multipart)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if err := c.writeParam(w, "from", c.from); err != nil {
		return err
	}
	if err := c.writeParam(w, "to", email); err != nil {
		return err
	}
	if err := c.writeParam(w, "subject", subject); err != nil {
		return err
	}
	if err := c.writeParam(w, "text", body); err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	// set header
	req.SetBasicAuth("api", c.apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return &MailgunError{
			Status: resp.StatusCode,
			Body:   buf.String(),
		}
	}
	return nil
}

func (c *client) writeParam(w *multipart.Writer, key, value string) error {
	fw, err := w.CreateFormField(key)
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(value))
	if err != nil {
		return err
	}
	return nil
}

func (e *MailgunError) Error() string {
	return fmt.Sprintf("Mailgun error[%d] %s", e.Status, e.Body)
}
