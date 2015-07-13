#carbonmax
a graphite client (string based carbon cache agent) for collecting and feeding data to carbon cache

##config file config.json
config.json should be placed in /etc/carbonmax/config.jsonor with -f flag
the example config.json is only for deployment of carbonmax through puppet

carbonlink: is for carbonmax(agent) setting

metric describe the metric and data which gernated by a system command or a sctipt
you can also create a script which will return the data you want to assign to the metric

please check syntax on config.json if you want multiple return!
OFS must be "|", and please match both part

##how to run or install
```
go build carbonmax.go
./carbonmax
```

###the agent can be run by cron or daemonlized by supervisord
####cron
```
add it to crontable
```
####supervisord
```
run carbonmax with config.json daemonize => true
```
this will run by the inerval you give in config.sjon

##future
if you have question, you can contact me at siegfried.chen@gmail.com


####The MIT License (MIT)

Copyright (c) 2014 Xuefeng Chen