# Shortcodes

Shortcodes allow template authors to provide extra functionality for content editors to display custom html/componenets within any description/synopsis fields on their site.

## Creation

Any jet templates that sit in the `/templates/shortcodes` directory of a site template are usable as a shortcode.

Each short code template can reference parameters provided to it as normal template variables. e.g:
`templates/shortcodes/code.go`
```html
<pre class="{{class}}">{{code}}</pre>
```

## Usage
From with a markdown editor of the admin site an editor can use a shortcode using the following format: `{{code class="js" code="alert('Hello World');"}}`

## Provided Shortcodes
The following shortcodes are provided by kibble.

### youtube
Embed a youtube video into a film (or other item) synopsis.

Example: `{{youtube id=dQw4w9WgXcQ class="never"}}`

### Parameters
<dl>
<dt>id</dt>
<dd>(Required) youtube video id.</dd>
<dt>class</dt>
<dd>(Optional) name of css class for styling (sizing/positioning) the embed.</dd>
</dl>

