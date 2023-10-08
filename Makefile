build:
	go build -o ./echo ./cmd/echo

test: build
	./maelstrom/maelstrom test -w echo --bin ./echo --nodes n1 --time-limit 10 --log-stderr
