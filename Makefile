.PHONY: all restaurant notification clean

all: restaurant notification

restaurant:
	cd restaurant && go mod tidy && go build

notification:
	cd notification && go mod tidy && go build

clean:
	rm -f restaurant/restaurant notification/notification
