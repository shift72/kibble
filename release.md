# Release

## Publish

Publish will zip all files placed in the ```/.kibble/dist``` directory

```bash

/.kibble
  /dist       - publish directory
  /kibble.zip - zip file to be published

```

## Releasing new versions - Shift72 internal only

Kibble is released to 3 places

  1. GitHub
  2. SHIFT72 Platform - this is where the platform will pull the kibble release from
  3. NPM - to support installation for third parties via NPM and Node.js environment

### a. Update version numbers and change log

```
changelog.md
kibble-npm/package.json
```

### b. Commit changes

```
git commit -a -m "release nnn"
```

### c. Tag Commit
  Ensure that the release is tagged correctly. Miss the prepended 'v' as this will mess S3 up.

```
git tag 0.9.6 master
git push origin 0.9.6
```

### d. build and release to locations 1 and 2

```
make release
```

### e. register new build in uber admin

Manual step, create rows for the staging and production versions
  * http://localhost:10001/admin/user~builder_version
  * http://localhost:10002/admin/user~builder_version


### f. update any sample templates with the new kibble version
