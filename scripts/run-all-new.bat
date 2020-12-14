@echo off
set ip=%~1
set /a firstIpPort=%~2
set /a lastIpPort=%~3
set /a firstUrlPort=%~4
set /a lastUrlPort=%~5
echo ip=%ip%
echo firstIpPort=%firstIpPort%
echo lastIpPort=%lastIpPort%
echo firstUrlPort=%firstUrlPort%
echo lastUrlPort=%lastUrlPort%
for /l %%i in (8001 1, 8006) do start run-new.bat -all %ip% %firstIpPort% %firstUrlPort%;
echo %firstIpPort%
echo %firstUrlPort%
set /a firstIpPort = %firstIpPort% + 1
set /a firstUrlPort = %firstUrlPort% + 1
pause