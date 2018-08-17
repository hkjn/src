# pki

Some simple tools for working with public key infrastructure tasks like certs, keys.

# Usage

## Generate self-signed Certificate Authority's certificate and key

```
docker run -v /etc/pki:/certs hkjn/pki initca
```

# Generate certificate and key for `myclient` signed by the CA

```
docker run -v /etc/pki:/certs hkjn/pki gencert myclient
```
