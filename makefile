source = linkdeep.go miner/common.go miner/config.go miner/fofa.go

all: dist

dist: dist/linkdeep_linux64 dist/linkdeep_linux32 dist/linkdeep_win64 dist/linkdeep_win32 dist/linkdeep_macOS

dist/linkdeep_linux64: $(source)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linkdeep_linux64

dist/linkdeep_linux32: $(source)
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o dist/linkdeep_linux32

dist/linkdeep_win64: $(source)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/linkdeep_win64

dist/linkdeep_win32: $(source)
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o dist/linkdeep_win32

dist/linkdeep_macOS: $(source)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/linkdeep_macOS

clean:
	rm -rf dist