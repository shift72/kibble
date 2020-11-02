# Change log

## 0.15.0
Stop trying to read entire S3 bucket contents on first sync (due to lack of index.kibble file).
Add `routeToPath(path string)` tempalte function for prepending current langauge (if not default) to relative path.

## 0.14.0
Add support for availability and pricing information

## 0.13.0
Skipped

## 0.12.0
Add support for pricing

## 0.11.1
Add 2 new functions to tv season and film crews, `GetJobNames` and `GetMembers` to get a list of unique crew job names, and a list of unique names of crew members with a particular job name

## 0.11.0
Add support for tags

## 0.10.22
Added support for duplicate film titles

## 0.10.21
Added support for Classifications in tv seasons model

## 0.10.20
Added support for `Studio` in the film model as a string array.

## 0.10.19
Fix goreleaser config for building version correctly.

## 0.10.18
Making sure version of kibble is set before publish.

## 0.10.17
Added support for markdown content allowing `<a />` tags to include `target="_blank"`.
Add support for markdown content to auto add `target="_blank"` to fully qualified links on `<a />` tags.

## 0.10.16
Added support for Classification records
Added field Classifications in film model, an array of Classification records.

## 0.10.15
Field `Subtitles` in film model is now an array. `GetSubtitles()` iterates through this and `SubtitleTracks`.

## 0.10.14
Move field Subtitles to field SubtitleTracks. Add new field Subtitles which captures the likes of hard-coded subtitles.
Add new function `GetSubtitles()` for tv episodes, films and bonus content. Populate this with a unique list of `SubtitleTracks.Name` and `Subtitles` values.

## 0.10.13
Import taglines and descriptions from bundles

## 0.10.12
Added support for design-time proxy routes. This allows a designer to set up patterns that will be managed by the reverse proxy instead of the static file matching. This allows the Player to be used locally (kind of).

## 0.10.11
Added support for standalone episode templates (full item and partial)

## 0.10.10
Added support for an SEO specific image.

## 0.10.9
Support Episode and ShowInfo Overviews being transformed for a season when its rendered.

## 0.10.8
Support season and film release_date that is a string in json, that fails to parse to a `time.Time`.

## 0.10.7
Changed live reload script to not rely on `If-Modified-Since` header, which seems to have issues in IE.
Added render option (`--serve`) to just serve site without a live reload script to ease testing (and not choke in IE).

## 0.10.6
Added a new GenericItem.GetTranslatedTitle function which allows specifying a i18n key

## 0.10.5
Vesion bump for a npm fix

## 0.10.4
Added type specific Get* methods for CustomFields

## 0.10.3
Added support for Bonus Content on TV Seasons
Added support via `map[string]interface{}` for Custom Fields (`custom` property) to Films, TV Season, Episodes, and Bonus Content.

## 0.10.2

Added support for config and configInt with defaults
Error logging includes file paths where possible
Added warning for i18n translations when a empty string is passed as a parameter
Validates datasource routes for valid replacement arguments
Validates templates for required and expected paths before publishing
Added command `kibble datasources` to print available datasources and paths

## 0.10.1

Added route validation

## 0.10.0

Add upstream remote to cloned repository
Removed --force option from init

## 0.9.11

Fix intermittent live reload bug when no status code is set

## 0.9.10

Windows fixes for template selection
Windows fixes for launching the browser

## 0.9.9

Updates to console messaging
Check builderVersion before publishing

## 0.9.8

Updates to the post install instructions

## 0.9.7

Season image fallbacks

## 0.9.6

Append base url to plans

## 0.9.5

Add Portrait Image to Plans

## 0.9.4

Add GetTitle to TvSeason
Support localising Tv Season Title
Removed pointers to optional Plan information

## 0.9.3

Add slug to TV Episode

## 0.9.2

Fix SEO support, based on changes made to:

* `/services/meta/v2/bios`
* `/services/meta/v1/bundles/`
* `/services/meta/v2/film/:ids/show_multiple`
* `/services/meta/v2/tv/seasons/show_multiple`

## 0.9.1

 Fix support for ints in i18n template function.

## 0.9.0

  Add support for plans - accessible via the Site

## 0.7.1

  Improve change detection speed

## 0.7.0

  Prevent syncing if there were errors durring rendering
  Added currentUrlPath

## 0.6.0

  Added Runtime.Localise

## 0.5.2

  Fixed resolving the collection items

## 0.2.1

  Moved short codes underneath the [site root]/templates/shortcodes
  Support subtitles
  Exclude path of '/' from zipfile
  TV Support
  Fixed bundle support
  Include errors in sync summary

## 0.2.0

  Fixed zip file creation, preserves structure
  Added zip only parameters for testing
  Support IAM profiles
