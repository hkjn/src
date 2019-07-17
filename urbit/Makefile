IMAGE_VERSION=0.1.0
.DEFAULT_GOAL=run
URBIT_DIR := "/usr/local/src/urbit"

build-image:
	docker build --tag hkjn/urbit:$(IMAGE_VERSION) .

clean:
	rm -rf $(HOME)/.urbit

init: build-image
	docker run --rm -it -v $(URBIT_DIR):/urbit hkjn/urbit:$(IMAGE_VERSION)

run: build-image
	docker run --rm -it -v $(URBIT_DIR):/urbit hkjn/urbit:$(IMAGE_VERSION) zod