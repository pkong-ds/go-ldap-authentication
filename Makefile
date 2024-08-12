.PHONY: tls
tls:
	mkcert localhost

.PHONY: copy-ca-cert
copy-ca-cert:
	cp "$$(mkcert -CAROOT)"/rootCA.pem .

include .env

# LDAP_URL=ldap://localhost:389
# LDAP_BASE_DN = cn=dataadmin,ou=datateam,dc=example,dc=org
# LDAP_BIND_DN = cn=admin,dc=example,dc=org
# LDAP_BIND_PASSWORD = admin
.PHONY: start
start:
	LDAP_URL=$(LDAP_URL) \
	LDAP_BASE_DN=$(LDAP_BASE_DN) \
	LDAP_BIND_DN=$(LDAP_BIND_DN) \
	LDAP_BIND_PASSWORD=$(LDAP_BIND_PASSWORD) \
	LDAP_SEARCH_QUERY=$(LDAP_SEARCH_QUERY) \
	go run main.go