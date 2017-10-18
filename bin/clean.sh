#!/bin/sh

for d in log state report ; do
	git rm -rf $d
	git rm -rf i686/${d}
	git rm -rf x86_64/${d}
	rm -rf $d
done 2> /dev/nul

for d in i686 x86_64 ; do
	git rm -rf ${d}/*.json
	rm -rf ${d}/*.json
done 2> /dev/nul

git rm -f ./*.exe
rm -f ./*.exe
