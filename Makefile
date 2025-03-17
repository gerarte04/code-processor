.PHONY: launch_services launch_services_with_tests stop_services build_services

launch_services: build_services
	docker compose up

launch_services_with_tests: build_services
	docker compose --profile test up --abort-on-container-exit --exit-code-from tester

stop_services:
	docker compose down

build_services:
	docker compose build
