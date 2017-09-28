#!/bin/sh

git rm -f log/*
git rm -f state/*
git rm -f report/*
git rm -f *.exe

rm -rf log state report *.exe
