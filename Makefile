all:
	cgogen pocketsphinx.yml

clean:
	rm -f pocketsphinx/cgo_helpers.go pocketsphinx/cgo_helpers.h pocketsphinx/doc.go pocketsphinx/types.go pocketsphinx/const.go
	rm -f pocketsphinx/pocketsphinx.go

test:
	cd pocketsphinx && go build

install:
	cd pocketsphinx && go install
	