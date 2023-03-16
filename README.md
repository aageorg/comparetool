# Compare tool for altlinux packages

## Dependencies
The tool uses [go-version](http://github.com/hashicorp/go-version) library for parsing and comparing package versions.
and [altclient](http://github.com/aageorg/altclient) for easy access to ALT Linux repository API

## Installation and Usage
Download package from github and build id:
```
$ git clone github.com/aageorg/altclient
$ cd comparetool
$ go build -o compare
```

Get json object with a differences between two ALT Linux branches:
```
./compare p10 p9
```
Layout is an array with two objects representing lists of missing and obsolete packages, grouped by architecture name:

```
[
  {
    "branch": "string",
    "missing": map["string"] [  // map where the key is the name of architecture, value - array of packages
     {                          // supported this architecture
      "name": "string",
      "epoch": 0,
      "version": "string",
      "release": "string",
      "arch": "string",
      "disttag": "string",
      "buildtime": 0,
      "source": "string"
      },
      ...
    ]
  },
  {
    "branch": "string",
    "missing": map["string"] [
     {
      "name": "string",
      "epoch": 0,
      "version": "string",
      "release": "string",
      "arch": "string",
      "disttag": "string",
      "buildtime": 0,
      "source": "string"
      },
      ...
    ]
    "out_of_date": map["string"] [
     {
      "name": "string",
      "epoch": 0,
      "version": "string",
      "release": "string",
      "arch": "string",
      "disttag": "string",
      "buildtime": 0,
      "source": "string"
      },
      ...
    ]
  }
]
```

## Issues and Contributing
If you find an issue with this library, please report an issue. If you'd like, we welcome any contributions. Fork this library and submit a pull request.