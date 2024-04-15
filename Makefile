.PHONY: test
test: 
	docker compose up --build --exit-code-from test

.PHONY: clean
clean:
	docker compose down

.PHONE: run-dependencies
run-dependencies:
	docker compose up bitcoin-core
