gencert-ca:
	@cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca

gencert-service:
	@cfssl gencert \
		-ca=$(TLS_DIR)/ca.pem \
		-ca-key=$(TLS_DIR)/ca-key.pem \
		-config=test/ca-config.json \
		-profile=$(service) \
		test/$(service)-csr.json | cfssljson -bare $(service)

gencert-client-apigw:
	@cfssl gencert \
		-ca=.tls/ca.pem \
		-ca-key=.tls/ca-key.pem \
		-config=test/ca-config.json \
		-profile=client-apigw \
		test/client-apigw-csr.json | cfssljson -bare client-apigw

test-curl:
	@curl -v https://localhost:8043/version \
		--cacert .tls/ca.pem \
		--key .tls/client-apigw-key.pem \
		--cert .tls/client-apigw.pem 