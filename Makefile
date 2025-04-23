build:
	templ generate
	go build -o ./tmp/main.exe .

run: build
	./tmp/main.exe

live:
	air -c .air.toml