.PHONY: test
test: 
	docker compose up --abort-on-container-exit --build 
