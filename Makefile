CERT_DIR=.ssl

gencert-ca:
	@cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca

gencert-service:
	@cfssl gencert \
		-ca=$(CERT_DIR)/ca.pem \
		-ca-key=$(CERT_DIR)/ca-key.pem \
		-config=test/ca-config.json \
		-profile=$(service) \
		test/$(service)-csr.json | cfssljson -bare $(service)

	@mv *.csr *.pem $(CERT_DIR)

gencert-client-apigw:
	@cfssl gencert \
		-ca=.tls/ca.pem \
		-ca-key=.tls/ca-key.pem \
		-config=test/ca-config.json \
		-profile=client-apigw \
		test/client-apigw-csr.json | cfssljson -bare client-apigw

test-curl:
	@curl -v https://localhost:8043/datetime?timestamp=1609417705 \
		--cacert $(CERT_DIR)/ca.pem \
		--key $(CERT_DIR)/client-apigw-key.pem \
		--cert $(CERT_DIR)/client-apigw.pem 