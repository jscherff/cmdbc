#!/bin/sh

git rm -f {log,state,report}/*

rm -rf ./{log,state,report}
rm -rf */{log,state,report,config.json}
