drone-dart
==========

Drone continuous delivery for Dart's Pub manager



## Storage Setup

Create a bucket:

```
gsutil mb -l US gs://BUCKET_NAME
```

Ensure the bucket is public read write:

```
gsutil -m setacl -R -a public-read gs://BUCKET_NAME
gsutil -m setdefacl public-read gs://BUCKET_NAME
```

Setup CORS by uploading the following `_cors.json` file to your bucket:

```
[
    {
      "origin": ["*"],
      "responseHeader": ["x-meta-goog-custom", "Access-Control-Allow-Origin"],
      "method": ["GET", "HEAD", "DELETE"],
      "maxAgeSeconds": 3600
    }
]
```

And then execute the following command:

```
gsutil cors set cors.json gs://BUCKET_NAME
```
