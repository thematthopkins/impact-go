openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout cert.pem \
    -new \
    -out cert.pem \
    -subj /CN=*.impact.dev \
    -reqexts SAN \
    -extensions SAN \
    -config <(cat /System/Library/OpenSSL/openssl.cnf \
        <(printf '[SAN]\nsubjectAltName=DNS:*.impact.dev')) \
    -sha256 \
    -days 3650
