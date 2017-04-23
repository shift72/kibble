# release new version

gox -osarch="darwin/amd64" -cgo -verbose -output="pkg/{{.OS}}_{{.Arch}}/kibble"

cd pkg

# zip files
for i in *
do
[ -d "$i" ] && zip -r "$i.zip" "$i"
done

# remove directories
rm -r */

# rename darwin - osx

# upload  - https://github.com/tcnksm/ghr