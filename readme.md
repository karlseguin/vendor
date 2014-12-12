# Go Package Vendoring

Vendor is meant to be a minimalist solution to managing dependencies in Go. It's similar to npm in that dependencies are stored within the project itself.

Vendor currently only works with git.

## Usage
First, add `.vendor` do you `.gitignore` file; vendored dependencies should not be committed.

Next, create a `vendor.json` file, which should be committed, in the root of your project. It has a simple format:

```json
{
  "typed": {
    "url": "https://github.com/karlseguin/typed.git",
    "revision": "60ea22ece11445c6ca01d44d414b84181144e072"
  },
  ...
}
```

After executing `vendor`, the dependencies will be available in your project's `.vendor` subfolder. You can import these using a relative path, or, preferably using something like: `example.com/yourproject/.vendor/typed`.
