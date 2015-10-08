package main

func main() {
	done := make(chan bool)

	go Httpd()
	go UdpConfig()

	<-done
}