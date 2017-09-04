all:
	c-for-go -maxmem 0x1fffffff pocketsphinx.yml

clean:
	rm -f pocketsphinx/cgo_helpers.go pocketsphinx/cgo_helpers.h pocketsphinx/cgo_helpers.c
	rm -f pocketsphinx/doc.go pocketsphinx/types.go pocketsphinx/const.go
	rm -f pocketsphinx/pocketsphinx.go

test:
	cd pocketsphinx && go build

install:
	cd pocketsphinx && go install
