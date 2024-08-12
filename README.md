# How to start this thing

Install mkcert

```
brew install mkcert
```

Run `make tls` to generate a TLS certificate and the key.

```
make tls
```

Run `make copy-ca-cert` to copy the mkcert root CA certificate in this project.

```
make copy-ca-cert
```

Run the LDAP server and the LDAP admin server.

```
docker compose up -d
```

For first time setup

1. Visit http://localhost:58080
2. Sign in with
  - DN: `cn=admin,dc=example,dc=org`
  - Password: `admin`
3. Create a POSIX group
4. Create as many users as you want.
