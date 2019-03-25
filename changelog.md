# Change log

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
