package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

type LDAPClient struct {
	conn *ldap.Conn

	host string
	port int
	bindDN string
	bindPassword string

	username string
	password string
}

func NewLDAPClient(host string, port int, bindDN string, bindPassword string, username string, password string) *LDAPClient {
	return &LDAPClient{
		host: host,
		port: port,
		bindDN: bindDN,
		bindPassword: bindPassword,

		username: username,
		password: password,
	}
}

func (l *LDAPClient) Connect() (error) {
	conn, err := ldap.DialURL(fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		log.Fatalf("Failed to connect to LDAP server: %s", err)
		return err
	}

	l.conn = conn

	// Bind to the LDAP server with read access
	err = l.conn.Bind(l.bindDN, l.bindPassword)
	if err != nil {
		log.Fatalf("Failed to bind to LDAP server: %s", err)
		return err
	}

	return nil
}

func (l *LDAPClient) StartTLSConnect() (error) {
	conn, err := ldap.DialURL(fmt.Sprintf("%s:%d", l.host, l.port))
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

	// Bind to the LDAP server with read access
	err = l.conn.Bind(l.bindDN, l.bindPassword)
	if err != nil {
		log.Fatalf("Failed to bind to LDAP server: %s", err)
		return err
	}

	return nil
}


func (l *LDAPClient) Close() {
	l.conn.Close()
}

func (l *LDAPClient) AuthenticateUser() (error) {
	// Search for the user
	searchRequest := ldap.NewSearchRequest(
		"cn=dataadmin,ou=datateam,dc=example,dc=org",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%v)", &ldap.AttributeTypeAndValue{Type: "uid", Value: l.username}),
		[]string{},
		nil,
	)

	sr, err := l.conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("User not found or too many entries returned")
	}

	userDN := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.conn.Bind(userDN, l.password)
	if err != nil {
		return fmt.Errorf("Failed to authenticate user: %s", err)
	}

	return nil
}

func main() {

	// Create a new LDAP client
	client := NewLDAPClient("ldap://localhost", 389, "cn=admin,dc=example,dc=org", "admin", "chrislee", "123")
	// client := NewLDAPClient("ldaps://localhost", 636, "cn=admin,dc=example,dc=org", "admin", "chrislee", "123")

	// Connect to the LDAP server
	err := client.StartTLSConnect()
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	} else {
		fmt.Println("Connected")
	}

	// Authenticate the user
	err = client.AuthenticateUser()
	if err != nil {
		log.Fatalf("Failed to authenticate user: %s", err)
	} else {
		fmt.Println("User Authenticated")
	}

	// Close the connection
	client.Close()
}