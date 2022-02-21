all: certs/ca-key.pem certs/server-key.pem certs/client-key.pem

#
# Certs
#

certs:
	mkdir -p certs

certs/ca-key.pem: certs
	cfssl gencert -initca ca.json | cfssljson -bare certs/ca

certs/server-key.pem: certs
	cfssl gencert \
		-ca=certs/ca.pem \
		-ca-key=certs/ca-key.pem \
		-config=ca-config.json \
		-profile=server \
		server.json | cfssljson -bare certs/server

certs/client-key.pem: certs
	cfssl gencert \
		-ca=certs/ca.pem \
		-ca-key=certs/ca-key.pem \
		-config=ca-config.json \
		-profile=client \
		client.json | cfssljson -bare certs/client

#
# Clean
#

clean-ca:
	rm -f certs/ca.*

clean-server:
	rm -f certs/server.*

clean-client:
	rm -f certs/client.*

clean: clean-ca clean-server clean-client
