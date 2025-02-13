NAME=coloride

all: build

.PHONY: clean run build

build:
	go build ./cmd/$(NAME)

run:
	go run ./cmd/$(NAME)

clean:
	@rm -f $(NAME) $(NAME).exe
