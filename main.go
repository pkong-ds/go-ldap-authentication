package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-ldap/ldap/v3"
)

type LDAPConfig struct {
	URL          *url.URL
	BaseDN       string
	BindDN       string
	BindPassword string
	SearchQuery  string
}

func NewLDAPConfigFromEnv() (*LDAPConfig, error) {
	urlStr := os.Getenv("LDAP_URL")
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &LDAPConfig{
		URL:          u,
		BaseDN:       os.Getenv("LDAP_BASE_DN"),
		BindDN:       os.Getenv("LDAP_BIND_DN"),
		BindPassword: os.Getenv("LDAP_BIND_PASSWORD"),
		SearchQuery:  os.Getenv("LDAP_SEARCH_QUERY"),
	}, nil
}

func (c *LDAPConfig) Connect() (*ldap.Conn, error) {
	return ldap.DialURL(c.URL.String())
}

func (c *LDAPConfig) StartTLS(conn *ldap.Conn) error {
	return conn.StartTLS(&tls.Config{
		ServerName: c.URL.Hostname(),
	})
}

func (c *LDAPConfig) Bind(conn *ldap.Conn) error {
	err := conn.Bind(c.BindDN, c.BindPassword)
	if err != nil {
		return fmt.Errorf("failed to authenticate search user: %w", err)
	}
	return nil
}

func (c *LDAPConfig) AuthenticateUser(conn *ldap.Conn, username string, password string) error {
	fmt.Println(c.SearchQuery)
	// Search for the user
	searchRequest := ldap.NewSearchRequest(
		c.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%v)", &ldap.AttributeTypeAndValue{Type: c.SearchQuery, Value: username}),
		[]string{},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("user not found or too many entries returned")
	}

	userDN := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = conn.Bind(userDN, password)
	if err != nil {
		return fmt.Errorf("failed to authenticate user: %s", err)
	}

	return nil
}

func doMain() (err error) {
	// get command line input for username and password
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	cfg, err := NewLDAPConfigFromEnv()
	if err != nil {
		err = fmt.Errorf("failed to create LDAP config from env: %w", err)
		return
	}

	var conn *ldap.Conn
	if cfg.URL.Scheme == "ldap" {
		conn, err = cfg.Connect()
		if err != nil {
			err = fmt.Errorf("failed to connect to LDAP server: %w", err)
			return
		}

		// Try StartTLS
		err = cfg.StartTLS(conn)
		if err != nil {
			log.Printf("failed to StartTLS: %v\n", err)
			log.Printf("fallback to plain LDAP\n")
			conn, err = cfg.Connect()
			if err != nil {
				err = fmt.Errorf("failed to connect to LDAP server: %w", err)
				return
			}
		} else {
			log.Printf("StartTLS\n")
		}
	} else if cfg.URL.Scheme == "ldaps" {
		conn, err = cfg.Connect()
		if err != nil {
			err = fmt.Errorf("failed to connect to LDAP server: %w", err)
			return
		}
	}
	defer conn.Close()
	log.Printf("connected\n")

	err = cfg.Bind(conn)
	if err != nil {
		return
	}
	log.Printf("bound\n")

	err = cfg.AuthenticateUser(conn, username, password)
	if err != nil {
		return
	}

	fmt.Println("user Authenticated")
	return nil
}

func main() {
	err := doMain()
	if err != nil {
		panic(err)
	}
}
