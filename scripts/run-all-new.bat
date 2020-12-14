@echo off
set ip=%~1
set /a firstIpPort=%~2
set /a lastIpPort=%~3
set /a firstUrlPort=%~4
set /a lastUrlPort=%~5
set /a usercnt=(lastIpPort - firstIpPort) + 1
echo usercnt=%usercnt%
echo ip=%ip%
echo firstIpPort=%firstIpPort%
echo lastIpPort=%lastIpPort%
echo firstUrlPort=%firstUrlPort%
echo lastUrlPort=%lastUrlPort%
set /a runRouter="true"
for /l %%i in (1, 1, %usercnt%) do start run-new.bat -all %ip% %firstIpPort% %firstUrlPort% %runRouter% user%%i;
echo %firstIpPort%
echo %firstUrlPort%
set /a firstIpPort = firstIpPort + 1
set /a firstUrlPort = firstUrlPort + 1
set /a runRouter="false"
pause