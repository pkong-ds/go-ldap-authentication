# Prerequisite

Making locally-trusted development certificates for localhost for local TLS connection.
This can be done by using `mkcert`

```
mkcert localhost
```

Then Mounting those file to the same directory.

ref: https://github.com/osixia/docker-openldap?tab=readme-ov-file#use-your-own-certificate

# Quick Start

To start the LDAP Server:

```
docker run --name ldap-server \
        --hostname ldap-server \
				--volume /PATH-TO-CERT:/container/service/slapd/assets/certs \
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
```

Then go to https://localhost:6443/ , login with

- Admin DN: `cn=admin,dc=example,dc=org`
- Admin Password: `admin`

Now you can config your own ldap server!

ref: https://blog.puckwang.com/posts/2022/use-docker-run-ldap-server/

To run the ldapclient, modify `.env` parameter for your ldap server configuration then run `make start`, then enter your login credentials in terminal.
