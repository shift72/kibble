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

You will be prompted to select a relevant version number.
Afterwards it'll ask you permission to do the different steps, you should say yes
to all these steps.

```
npx release-it
```

### c. Push new site builder version to db

After a new kibble build has been released, it needs to be added to the versions
available on the site builder for both staging and prod.

The publish release action automatically pushes it to staging, but the prod
release is still manual.

The easy way is the Github action: "Publish Site Builder Version".

Pick the env and enter the kibble semver version e.g. "0.17.7" (no v prefix!)

OR, do it manually via uberadmin:

Manual step, create rows for the staging and production versions
  * http://localhost:10001/admin/user~builder_version
  * http://localhost:10002/admin/user~builder_version


### d. Update any sample templates with the new kibble version
