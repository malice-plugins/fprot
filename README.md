malice-fprot
============

[![Circle CI](https://circleci.com/gh/maliceio/malice-fprot.png?style=shield)](https://circleci.com/gh/maliceio/malice-fprot) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/fprot.svg)](https://hub.docker.com/r/malice/fprot/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/fprot.svg)](https://hub.docker.com/r/malice/fprot/) [![Docker Image](https://img.shields.io/badge/docker%20image-299%20MB-blue.svg)](https://hub.docker.com/r/malice/fprot/)

This repository contains a **Dockerfile** of [fprot](http://www.fprot.net/lang/en/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/fprot/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

-	[ubuntu:xenial (*118 MB*\)](https://hub.docker.com/r/library/ubuntu/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/fprot/) from public [DockerHub](https://hub.docker.com): `docker pull malice/fprot`

### Usage

```
docker run --rm malice/fprot EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/fprot FILE

Usage: fprot [OPTIONS] COMMAND [arg...]

Malice F-PROT AntiVirus Plugin

Version: v0.1.0, BuildTime: 20170122

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V         verbose output
  --table, -t	        output as Markdown table
  --callback, -c	    POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x	        proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --timeout value       malice plugin timeout (in seconds) (default: 60) [$MALICE_TIMEOUT]    
  --elasitcsearch value elasitcsearch address for Malice to store results [$MALICE_ELASTICSEARCH]   
  --help, -h	        show help
  --version, -v	        print the version

Commands:
  update	Update virus definitions
  web       Create a fprot scan web service  
  help		Shows a list of commands or help for one command

Run 'fprot COMMAND --help' for more information on a command.
```

Sample Output
-------------

### JSON:

```json
{
  "f-prot": {
    "infected": true,
    "result": "EICAR_Test_File (exact)",
    "engine": "4.6.5.141",
    "updated": "20170122"
  }
}
```

### STDOUT (Markdown Table):

---

#### F-PROT

| Infected | Result                  | Engine    | Updated  |
|----------|-------------------------|-----------|----------|
| true     | EICAR_Test_File (exact) | 4.6.5.141 | 20170122 |

---

Documentation
-------------

-	[To write results to ElasticSearch](https://github.com/maliceio/malice-fprot/blob/master/docs/elasticsearch.md)
-	[To create a F-PROT scan micro-service](https://github.com/maliceio/malice-fprot/blob/master/docs/web.md)
-	[To post results to a webhook](https://github.com/maliceio/malice-fprot/blob/master/docs/callback.md)
-	[To update the AV definitions](https://github.com/maliceio/malice-fprot/blob/master/docs/update.md)

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-fprot/issues/new).

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice-fprot/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice-fprot/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice-fprot/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

MIT Copyright (c) 2016-2017 **blacktop**
