package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// Send mail ( by bt )

func SendMailRaw(addr string, port string, heloDomain string, smtpAuth smtp.Auth, from string, to string, emailMsg *[]byte) error {
	if len(*emailMsg) == 0 {
		return fmt.Errorf("msg is nil")
	}
	if len(port) != 0 && port[0:1] != ":" {
		port = ":" + port
	}

	// 2. dial
	client, err := smtp.Dial(addr + port)
	if err != nil {
		return err
	}
	defer func(client *smtp.Client) {
		err := client.Close()
		if err != nil {
			log.Println("failed to close connection")
		}
	}(client)

	// 3. cert
	{
		// hello
		if heloDomain != "" {
			if err = client.Hello(heloDomain); err != nil {
				return err
			}
		}

		// start tls
		isExistTLS, _ := client.Extension("STARTTLS")
		if isExistTLS {
			config := &tls.Config{
				ServerName:         addr,
				InsecureSkipVerify: true,
			}
			if err = client.StartTLS(config); err != nil {
				return err
			}
		}

		// auth
		if smtpAuth != nil {
			isExistExtAuth, _ := client.Extension("AUTH")
			if !isExistExtAuth {
				return fmt.Errorf("smtp server doesn't support AUTH")
			}
			if err = client.Auth(smtpAuth); err != nil {
				return err
			}
		}
	}

	// 4. Send
	{
		if err := client.Mail(from); err != nil {
			return err
		}

		if err = client.Rcpt(to); err != nil {
			return err
		}

		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(*emailMsg)

		if err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}
	}
	return nil
}
