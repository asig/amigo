all: amigo.exe amigo_x64.exe

amigo.exe: amigo.go
	GOOS=windows GOARCH=386 go build -o amigo.exe

amigo_x64.exe: amigo.go
	GOOS=windows GOARCH=amd64 go build -o amigo_x64.exe