# cnodes

cnodes is command line tool to print consul nodes with health status.

```
$ cnodes
critical  172.18.0.3       client_1 dc1
passing   172.18.0.4       client_2 dc1
passing   172.18.0.2       server dc1
```

## Config

```
$ cat ~/.cnodes
[DEFAULT]
url = http://localhost:8500/

[dc_1]
url = http://dc.1:8500/

[dc_2]
url = http://dc.2:8500/

[dc_3]
url = http://dc.3:8500/
```
