A hello world project powered by [NES.css] and [Netlify Functions].

[NES.css]: https://nostalgic-css.github.io/NES.css/
[Netlify Functions]: https://functions.netlify.com/

The frontend part is a simple HTML page showing a progress bar.

The input is recerived from the hash part of the URL:

```
https://progressbar.mmap.page/#{%22url%22:%20%22https://example.com/api%22,%20%22value_key_name%22:%20%22foo%22,%20%22max_key_name%22:%22bar%22}
```

Here progressbar.mmap.page is where this project is hosted.

`https://example.com/api` should be a site return something like:

```json
{
    "foo": 123,
    "bar": 456,
    "other_keys": "will be ignored",
}
```

And the value of `foo` and `bar` will be used as the `value` and `max` attribute of the progress bar element (`<progress>`).

On the frontend part, it uses [NES.css] for styling and a few lines of JavaScript (without any framework or library) to parse user input and invoke the Netlify Function.

The Netlify Function is written in Go, which fetches the data from the requested URL.

