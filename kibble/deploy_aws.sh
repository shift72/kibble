VERSION=$1
OS=$2

if [ -z $VERSION ]; then
  echo "error: version not specified"
  exit 1
fi

if [ -z $OS ]; then
  echo "error: OS not specified"
  exit 1
fi

if [ $OS != "linux" ]; then
  echo "Skipping hook for $OS"
  exit 0
fi

echo "uploading linux only version - $VERSION"
aws s3 cp ./dist/kibble_linux_amd64_v1/kibble s3://shift72-sites/builder/$VERSION/kibble --profile shift72a