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

test-curl:
	@curl -v https://localhost:8081/datetime?timestamp=1609417705 \
		--cacert $(CERT_DIR)/ca.pem \
		--key $(CERT_DIR)/service-apigw-key.pem \
		--cert $(CERT_DIR)/service-apigw.pem 

docker-compose-up:
	@docker-compose -f docker-compose.yml up -V \
  		--build \
  		--abort-on-container-exit \
  		--remove-orphans