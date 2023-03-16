# Compare tool for ALT Linux packages

## Dependencies
The tool uses [go-version](http://github.com/hashicorp/go-version) library for parsing and comparing package versions
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
Example:
```
$ ./compare p10 p9 | more
Downloading of package list from branch p9 is started
Downloading of package list from branch p10 is started
Packages list from branch p10 is downloaded
Packages list from branch p9 is downloaded
Search of missing and obsolete packages in branch p9 is started
Search of missing packages in branch p10 is finished
Search of missing and obsolete packages in branch p9 is finished
[{"branch":"p10","missing":{"aarch64":[{"name":"python3-module-utmp","epoch":0,"version":"0.8","release":"alt1.1.1","arch":"aarch64","disttag":"sisyphus+2256
25.17300.91.1","buildtime":1555309208,"source":"python-module-utmp"},{"name":"cairo-dock-GMenu","epoch":0,"version":"3.4.1","release":"alt13","arch":"aarch64
","disttag":"sisyphus+228347.1400.1.2","buildtime":1556568828,"source":"cairo-dock-plugins"}...
```

Use flag `-q` if you want to get output without additional messages

## Issues and Contributing
If you find an issue with this library, please report an issue. If you'd like, we welcome any contributions. Fork this library and submit a pull request.