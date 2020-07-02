VERSION=$(git describe --tags)
if [ -z $VERSION ]; then
  echo "error: tagged version not found"
  exit 1
fi

echo "uploading linux only version - $VERSION"
aws s3 cp ./dist/shift72-kibble_linux_amd64/kibble s3://shift72-sites/builder/$VERSION/kibble --profile shift72a