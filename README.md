#carbonmax
a graphite client (string based carbon cache agent) for collecting and feeding data to carbon cache

##config file carbonmax.ini
carbonmax.ini should be placed in /etc/carbonmax.ini or with -inifile flag
the example carbonmax.ini is only for deployment of carbonmax through puppet

[carbonlink] is for carbonmax(agent) setting

[resouces] describe the mertic and data which gernated by a system command or a sctipt
you can also create a script which will return the data you want to assign to the mertic


##how to run or install
```
go run carbonmax.go -loop
go build carbonmax.go
```

###the agent can be run by cron or daemonlized by supervisord
####cron
```
add it to crontable
```
####supervisord
```
go run carbonmax -loop
```
this will run by the inerval you give in inifile

##future
???
if you have question, you can contact me at siegfried.chen@gmail.com


####The MIT License (MIT)

Copyright (c) 2014 Xuefeng Chen