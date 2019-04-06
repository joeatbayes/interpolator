set GOPATH=%cd%

go get -u "github.com/joeatbayes/goutil/jutil"



$ for GOOS in darwin linux windows solaris; do
    for GOARCH in 386 amd64; do
        go build interpolate/interpolate.go -v -o interpolate-$GOOS-$GOARCH
    done
done

$ for GOOS in solaris; do
    for GOARCH in sparc sparc64 386 amd64; do
        go build interpolate/interpolate.go -v -o interpolate-$GOOS-$GOARCH
    done
done
