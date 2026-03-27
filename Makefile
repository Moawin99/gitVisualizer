build:
	go build .

run: build
	./gitVisualizer $(ARGS)

clean:
	rm gitVisualizer