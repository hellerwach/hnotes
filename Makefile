default: build compress

build:
	go build -a -tags "netgo,osusergo" -ldflags='-w -s -extldflags "-static"' -gcflags=all="-l -B" \
	-o hnotes main.go
	
compress:
	upx -5 -k hnotes
	# I do not why, but upx sometimes creates hnotes.~
	rm -f hnotes.\~
	
clean:
	rm -f hnotes

readme:
	embedmd -w README.md