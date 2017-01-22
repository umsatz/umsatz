# umsatz frontend

Single Page Frontend for Umsatz. Written using Backbone Marionette.

## Run the tests

```
karma start
```

## Uploads & Downloads

curl -s -X POST --data-binary @fixture/test.zip 'http://umsatz.dev/upload/2014/foo.zip'
curl -s 'http://umsatz.dev/download/2014/foo.zip' > foo.zip