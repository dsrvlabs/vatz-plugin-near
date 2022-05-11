build:
	cd protocol-status-plugin
	go build main.go
	cd ..
	cd server-status-plugin
	go build main.go

clean:
	rm vatz
