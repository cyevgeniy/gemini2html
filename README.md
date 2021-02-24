# gemini2html
Convert gemini files to html.

## Usage

```
go run gemini2html.go
```

It will expect ```posts``` directory with *.gmi or *.gemini files and ```assets``` directory.

Generated site will lie in ```_site``` directory.

All content from ```assets``` directory will be copied to ```_site/assets/``` directory.
Genini2html expects that assets content is flat and has not nested directories. Same rule
for ```posts``` directory. All *.gmi files that lies in root of working directory, will be
converted and placed into root of _site directory.
