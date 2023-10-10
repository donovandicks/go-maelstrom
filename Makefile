.PHONY: echo broadcast

build-echo:
	go build -o ./echo ./cmd/echo

echo: build-echo
	./maelstrom/maelstrom test -w echo --bin ./echo --nodes n1 --time-limit 10 --log-stderr

build-broadcast:
	go build -o ./broadcast ./cmd/broadcast

broadcast: build-broadcast
	./maelstrom/maelstrom test -w broadcast --bin ./broadcast --time-limit 5 --log-stderr

broadcast-high: build-broadcast
	./maelstrom/maelstrom test -w broadcast --bin ./broadcast --time-limit 20 --rate 100
