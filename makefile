PHONY: build run

build:
	go build -o whatsappWebAPI .

run:
	make build
	./whatsappWebAPI
