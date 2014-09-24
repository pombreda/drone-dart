drone-dart
==========

Drone continuous delivery for Dart's Pub manager

## Runtime

Install the following runtime requirements:

* libsqlite3-dev
* docker
* mysql (optional)

## Building

Use the following commands to build from source:

```sh
make deps
make
```

To compile style sheets (requires less):

```sh
make lessc
```

To bundle static files inside the binary:

```sh
make embed
```

## Running

Run with local configuration:

```sh
./drone-dart
```

Customize the Basic Auth settings used to restict access to the system.
The default username / password is `admin:admin` primarily intended for
local testing. To customize:

```sh
./drone-dart --password="user:password"
```

Customize the database. The default database is SQLite3 primarily intended
for local testing. To customize or use an alternate database like MySQL:

```sh
./drone-dart --driver="mysql" --datastore="user:password@/dbname"
``` 
