@echo off
set ip=%~1
set ipPort=%~2
set urlPort=%~3
rem router param is false by default, can be changed to true if no router yet
set router=%~4
set username=%~5
start go run ../../raft-db/cmd/main.go -all %ip% %ipPort% %urlPort% %router% %username%
exit