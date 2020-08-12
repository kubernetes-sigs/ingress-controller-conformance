# Report builder for conformance tests

### Environment variables:

#### Mandatory

| Variable  | Description |
| ------------- | ------------- |
| BUILD  | String with identification about the build. Used to add information about the trend |
| INPUT_DIRECTORY  | Directory that contains the cucumber json files  |
| OUTPUT_DIRECTORY | Directory where the reports will be generated |

#### Optional

| Variable  | Description |
| ------------- | ------------- |
| RELEASE  | Information about the release where the tests where executed |



## Building

```console
make
```

### Generation of reports

```console
docker run \
    -e BUILD=$(git rev-parse --short HEAD) \
    -e INPUT_DIRECTORY=/input \
    -e OUTPUT_DIRECTORY=/output \
    -v $PWD:/input:ro \
    -v $PWD/output:/output \
    local/reports-builder:0.0
```

### Display

The reports are plain HTML files. The file located in `OUTPUT_DIRECTORY/index.html` renders the initial page.

Using any web server capable of render html is enough.
Like:

```console
cd $OUTPUT_DIRECTORY

python -m http.server 8000
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...

```
