set -e
echo "building"
GOOS=linux GOARCH=amd64 go build
chmod +x ./kibble

VERSION=$(grep -o '\d*\.\d*\.\d*' version/version.go)
echo "uploading - $VERSION"

aws s3 cp ./kibble s3://shift72-sites/builder/$VERSION/kibble --profile shift72a

