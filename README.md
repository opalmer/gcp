Disclaimer
----------
This is an experimental tool which was mainly intended to 
be a learning exercise in Go.  It was designed for a particular
purpose, producing files for backup.  Contributions or other 
improvements are still welcome however.

GCP
---

[![Build Status](https://travis-ci.org/opalmer/gcp.svg)](https://travis-ci.org/opalmer/gcp.svg)


gcp is a command line tool designed to implement features similar 
to the cp command on Linux in some respects. 

The features of gcp are:

  * Path exclusion or inclusion can be defined in a configuration file.
  * Multi-threading support.
  * Encryption (AES) and compression (lzma) of files on the fly.
  
  

Configuration
-------------

gcp looks for a configuration file in one of two places after loading 
the default:

```ini
[gcp]
encrypt = true
compress = true
dry_run = false
crypto_key =
include =
exclude = .DS_Store,.git,.svn,.hg,.egg*,__pycache__,.idea,*.pyc
```

By default gcp will load the default, then the file present in the 
``GCP_CONFIG`` environment variable then any file provided to ``-config`` on
the command line.