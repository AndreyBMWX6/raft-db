@echo off

set ip=%~1
set ipPort=%~2
set urlPort=%~3
set username=%~4
start go run ../../raft-db/cmd/main.go %ip% %ipPort% %urlPort% %username%
exit
