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

### a. Document changes

Within the changelog should be an [Unreleased] section 
add this or add your changes to the current one

```
changelog.md
```

### b. Run release-it

```
// You will be prompted to select a relevant version number
// Afterwards it'll ask you permission to do the different steps, for now only say no to 
// creating a release as it's currently broken due to the changelog being massive.

npx release-it
```

### c. Register new build in uber admin

Manual step, create rows for the staging and production versions
  * http://localhost:10001/admin/user~builder_version
  * http://localhost:10002/admin/user~builder_version


### d. Update any sample templates with the new kibble version
