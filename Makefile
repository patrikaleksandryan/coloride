NAME=coloride

all: build

.PHONY: clean run build

build:
	go build ./cmd/$(NAME)

run:
	GODEBUG=asyncpreemptoff=1 go run ./cmd/$(NAME)

clean:
	@rm -f $(NAME) $(NAME).exe
