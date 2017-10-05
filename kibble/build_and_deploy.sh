set -e
echo "building"
GOOS=linux GOARCH=amd64 go build
chmod +x ./kibble
echo "uploading"
aws s3 cp ./kibble s3://shift72-sites/builder/0.2.0/kibble --profile shift72a


#aws s3 cp s3://shift72-sites/builder/0.2.0/kibble . 
#sudo chmod +x ./kibble

