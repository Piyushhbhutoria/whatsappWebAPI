package main

type SendText struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
}

type SendImage struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
	Image    string `json:"image"`
}
