VERSION    := $(shell git describe --tags)

DARWINx64   := "s72-web/kibble/$(VERSION)/kibble_$(VERSION)_macOS_64-bit.zip"
DARWINarm64 := "s72-web/kibble/$(VERSION)/kibble_$(VERSION)_macOS_arm64-bit.zip"
LINUXx64    := "s72-web/kibble/$(VERSION)/kibble_$(VERSION)_Tux_64-bit.tar.gz"
LINUXarm64 := "s72-web/kibble/$(VERSION)/kibble_$(VERSION)_Tux_arm64-bit.tar.gz"
WINDOWSx64  := "s72-web/kibble/$(VERSION)/kibble_$(VERSION)_windows_64-bit.zip"
WINDOWSarm64  := "s72-web/kibble/$(VERSION)/kibble_$(VERSION)_windows_arm64-bit.zip"

release_with_profile: AWS_PROFILE=shift72a release

release:
	cd kibble && goreleaser --rm-dist

	echo "setting acls for the released versions"
	aws s3api put-object-acl --bucket shift72-sites --key $(DARWINx64)  --acl public-read
	aws s3api put-object-acl --bucket shift72-sites --key $(DARWINarm64)  --acl public-read
	aws s3api put-object-acl --bucket shift72-sites --key $(LINUXx64)   --acl public-read
	aws s3api put-object-acl --bucket shift72-sites --key $(LINUXarm64)   --acl public-read
	aws s3api put-object-acl --bucket shift72-sites --key $(WINDOWSx64) --acl public-read
	aws s3api put-object-acl --bucket shift72-sites --key $(WINDOWSarm64) --acl public-read
	aws s3 cp ./kibble/dist/kibble_linux_amd64_v1/kibble s3://shift72-sites/builder/$(VERSION)/kibble

update_s3:
	echo "setting acls for the released versions"
	aws s3api put-object-acl --bucket shift72-sites --key $(DARWINx64)  --acl public-read --profile shift72a
	aws s3api put-object-acl --bucket shift72-sites --key $(DARWINarm64)  --acl public-read --profile shift72a
	aws s3api put-object-acl --bucket shift72-sites --key $(LINUXx64)   --acl public-read --profile shift72a
	aws s3api put-object-acl --bucket shift72-sites --key $(WINDOWSx64) --acl public-read --profile shift72a
	