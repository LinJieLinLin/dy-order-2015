@echo off

cd %cd%
set GOPATH=%GOPATH%;%cd%
dogo
