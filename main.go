package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
)

type LDAPClient struct {
	conn *ldap.Conn

	URL          string
	baseDN       string
	bindDN       string
	bindPassword string
	searchQuery  string
}

func NewLDAPClient() *LDAPClient {
	return &LDAPClient{
		URL:          os.Getenv("LDAP_URL"),
		baseDN:       os.Getenv("LDAP_BASE_DN"),
		bindDN:       os.Getenv("LDAP_BIND_DN"),
		bindPassword: os.Getenv("LDAP_BIND_PASSWORD"),
		searchQuery:  os.Getenv("LDAP_SEARCH_QUERY"),
	}
}

func (l *LDAPClient) Connect() error {
	conn, err := ldap.DialURL(l.URL)
	if err != nil {
		log.Fatalf("Failed to connect to LDAP server: %s", err)
		return err
	}

	l.conn = conn
	return nil
}

func (l *LDAPClient) Bind() error {
	err := l.conn.Bind(l.bindDN, l.bindPassword)
	if err != nil {
		log.Fatalf("Failed to bind to LDAP server: %s", err)
		return err
	}

	return nil
}

func (l *LDAPClient) StartTLSConnect() error {
	conn, err := ldap.DialURL(l.URL)
	if err != nil {
		log.Fatalf("Failed to connect to LDAP server: %s", err)
		return err
	}

	l.conn = conn

	err = l.conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatalf("Failed to start TLS connection: %s", err)
		return err
	}

	return nil
}

func (l *LDAPClient) Close() {
	l.conn.Close()
}

func (l *LDAPClient) AuthenticateUser(username string, password string) error {
	fmt.Println(l.searchQuery)
	// Search for the user
	searchRequest := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%v)", &ldap.AttributeTypeAndValue{Type: l.searchQuery, Value: username}),
		[]string{},
		nil,
	)

	sr, err := l.conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("user not found or too many entries returned")
	}

	userDN := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.conn.Bind(userDN, password)
	if err != nil {
		return fmt.Errorf("failed to authenticate user: %s", err)
	}

	return nil
}

func main() {
	// get command line input for username and password
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	var err error
	// Create a new LDAP client
	client := NewLDAPClient()
	// Check if url starts with ldap://
	if client.URL[:7] == "ldap://" {
		// Connect to the LDAP server
		err = client.StartTLSConnect()
		if err != nil {
			log.Printf("Failed to connect: %s", err)
			err = client.Connect()
			if err != nil {
				log.Fatalf("Failed to connect: %s", err)
			} else {
				fmt.Println("Connected")
			}
		}
	} else if client.URL[:8] == "ldaps://" {
		// Connect to the LDAP server
		err = client.Connect()
		if err != nil {
			log.Fatalf("Failed to connect: %s", err)
		} else {
			fmt.Println("Connected")
		}
	}

	// Bind to the LDAP server
	err = client.Bind()
	if err != nil {
		log.Fatalf("Failed to bind: %s", err)
	} else {
		fmt.Println("Bound")
	}

	// Authenticate the user
	err = client.AuthenticateUser(username, password)
	if err != nil {
		log.Fatalf("Failed to authenticate user: %s", err)
	} else {
		fmt.Println("User Authenticated")
	}

	// Close the connection
	client.Close()
}

/*
docker run --name ldap-server \
        --hostname ldap-server \
				--volume /Users/york/Desktop/Oursky/certs:/container/service/slapd/assets/certs \
				--env LDAP_TLS_CRT_FILENAME=my-ldap.crt \
				--env LDAP_TLS_KEY_FILENAME=my-ldap.key \
				--env LDAP_TLS_CA_CRT_FILENAME=the-ca.crt \
				--env LDAP_TLS_VERIFY_CLIENT="try" \
		-p 389:389 -p 636:636 \
		--detach \
		osixia/openldap:latest

docker run --name ldap-admin \
    -p 6443:443 \
    --link ldap-server:ldap-host \
    --env PHPLDAPADMIN_LDAP_HOSTS=ldap-host \
    --detach \
    osixia/phpldapadmin:latest
*/
