# Go Package Vendoring

You probably shouldn't use this.

## Usage
First, add `.vendor` to you `.gitignore` file; vendored dependencies should not be committed.

Next, create a `vendor.json` file, which should be committed, in the root of your project. It has a simple format:

```json
{
  ".": "IMPORT OF CURRENT PROJECT",
  "typed": {
    "url": "https://github.com/karlseguin/typed.git",
    "revision": "60ea22ece11445c6ca01d44d414b84181144e072"
  },
  "bytepool": {
    "url": "https://github.com/karlseguin/bytepool.git",
    "revision": "d858cd4db848fa8f5275b746d57e84571a5c0be1"
  }
}
```

After executing `vendor`, the dependencies will be available in your project's `.vendor` subfolder. You can import these using a relative path, or, preferably using something like: `example.com/yourproject/.vendor/typed`.

## Sublime Find
In sublime, you can have .vendor automatically ignored when doing a search by adding the following to your preferences:

```json
"folder_exclude_patterns": [".vendor"]
```
