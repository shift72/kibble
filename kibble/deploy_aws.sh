VERSION=$(git describe --tags)
echo "uploading - $VERSION"
aws s3 cp ./dist/linuxamd64/kibble s3://shift72-sites/builder/$VERSION/kibble --profile shift72a