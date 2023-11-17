testing:
	./coverage.sh -m atomic testing

docker:
	docker build -f ./build/dockerfile -t tinder_match_system .

doc:
	sh swag.sh