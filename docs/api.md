# API

All templates are rendered using [JET](https://github.com/CloudyKit/jet/).

Any builtin function from JET is also available in kibble templates. The [JET Wiki](https://github.com/CloudyKit/jet/wiki) provides some helpfule examples and lists of usable functions.

## Template functions

<dl>
<dt>i18n(*key*, _default | pluralCount | templateMap_)</dt>
<dd>Prints translated text for a given key. If a value can't be found then an error is returned.

Basic Usage: `i18n("some_key")`

Optionally, an extra parameter can be provided to help template editors. The data type of this parameter makes the `i18n` act differently.

### default
If the parameter is string (default), then this value is used if a value could not be found for the provided key. e.g. `{{i18n("nonexistent_key", "This works")}}`

### pluralCount
If the parameter is a number, then this value is used to select the right plural form from the translation file. e.g. `{{i18n("library_item_count", 3)}}`

### templateMap
If the parameter is an object/hash, then the values are used by the translation template. e.g. `{{i18n("lots_of_things", map("Name", "Graham", "Country", "NZ"))}}` will print `I am Graham from NZ` with this translation string `"I am {{.Name}} from {{.Country}}"`.
</dd>

<dt>version</dt>
<dd>Prints the version of Kibble that rendered the current template.

Usage: `{{version}}`</dd>
<dt>lang</dt>
<dd>Prints the current language name

Usage: `{{lang}}`</dd>
<dt>routeToSlug(*slug*)</dt>
<dd>Prints the url path for a given items (film, bundle, episode, etc) slug.

Usage: `<a href="{{routeToSlug(film.Slug)}}">{{film.Title}}</a>`</dd>
<dt>routeTo(*item*)</dt>
<dd>Prints the url path to a specific item (film, bundle, episode, etc) page. Will print an error if the path could not be found.

Usage: `<a href="{{routeTo(film)}}">{{film.Title}}</a>`</dd>

<dt>config(*key*, _defaultValue_)</dt>
<dd>Prints the string value for the given config key. An optional default value can be provided if the key does not exist.</dd>
<dt>configInt(*key*, _defaultValue_)</dt>
<dd>Prints the number value for the given config key. An optional default value can be provided if the key does not exist.</dd>
<dt>isEnabled(*key*)</dt>
<dd>Prints whther the specific feature toggle value is enabled or not.</dd>
</dl>

## Models

### Site

### Film
Holds information for a given film.

#### Properties
| Name            | Type                   | Description |
|:----------------|:-----------------------|:------|
| ID              | int                    | |
| Slug            | string                 | |
| Title           | string                 | |
| TitleSlug       | string                 | Kebab case version of the film title, used for creating urls to the films detail page. |
| Trailers        | []Trailer              | |
| Bonuses         | BonusContentCollection | |
| Cast            | []CastMember           | |
| Crew            | CrewMembers            | |
| Studio          | []string               | |
| Overview        | string                 | This field is Markdown and can make use of shortcodes. |
| Tagline         | string                 | |
| ReleaseDate     | time.Time              | |
| Runtime         | Runtime                | |
| Countries       | StringCollection       | |
| Languages       | StringCollection       | |
| Genres          | StringCollection       | |
| Tags            | StringCollection       | |
| Seo             | Seo                    | |
| Images          | ImageSet               | |
| Prices          | PriceInfo              | |
| Available       | Period                 | |
| Recommendations | []GenericItem          | |
| Subtitles       | []string               | |
| SubtitleTracks  | []SubtitleTrack        | |
| CustomFields    | CustomFields           | |
| Classifications | []Classification       | |

#### Methods

