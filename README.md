# whatsappWebAPI

This project is made with the help of [Rhymen/go-whatsapp - WhatsApp Web API](https://github.com/Rhymen/go-whatsapp).
Works only for numbers in India for now.

## Requirements

1. Golang setup on your local
2. Web Browser/Postman

## Setup

1. Download the latest binary release.  
2. run on your system.  
3. Create a files folder in the directory where binary is running

or create your binary locally as below

1. Download the repo
2. run 
```shell
go get .
go run main.go
````

## Steps to use - 

1. Scan the QR code with whatsapp web
2. Put your lists for bulk message and pictures to send in files folder
3. Use the Following links to send text/image

## Usage

Hit the following links to send message on whatsapp  
Send Text --> localhost:8081/sendText?to=1234567890&msg=hello  
Send Image --> localhost:8081/sendImage?to=1234567890&msg=hello&img=testImg.jpg  
Send bulk text --> localhost:8081/sendBulk?file=test.csv  
Send bulk image --> localhost:8081/sendBulkImg?file=testImg.csv  

Link to import Postman requests --> https://www.getpostman.com/collections/caee3c3c8a3fc04304d0

### Demo

Test Text --> [localhost:8081/sendText](localhost:8081/sendText)  
Test Image --> [localhost:8081/sendText](localhost:8081/sendText) 

Demo bulk file examples stored in [files](/files) folder

#### Extra

Sample with Sentry and bot integrated in hosting folder
