#!/bin/bash

grep -rl "github.com/astaxie/beego" . --exclude-dir=.svn | xargs sed -i 's/github.com\/astaxie\/beego/github.com\/guerillagrow\/beego/g'
