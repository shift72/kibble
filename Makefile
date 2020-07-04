release:
	cd kibble && AWS_PROFILE=shift72a goreleaser --rm-dist
	cd ../kibble-npm && npm publish
