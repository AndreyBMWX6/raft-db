@echo off

set ip=%~1
set firstIpPort=%~2
set lastIpPort=%~3
set firstUrlPort=%~4
set lastUrlPort=%~5
set /a usercnt=(%lastIpPort% - %firstIpPort%) + 1
set /a usernum=1

echo running %usercnt% servers on ip:%ip%
echo ip=%ip%
echo firstIpPort=%firstIpPort%
echo lastIpPort=%lastIpPort%
echo firstUrlPort=%firstUrlPort%
echo lastUrlPort=%lastUrlPort%

start go run ../../internal/router/run/run_router.go

for /l %%i in (1, 1, %usercnt%) do (
call :run
call :inc
)
goto :eof

:inc
set /a firstIpPort=firstIpPort + 1
set /a firstURLPort=firstURLPort + 1
set /a usernum=usernum + 1
goto :eof

:run
start run.bat %ip% %firstIpPort% %firstUrlPort% user%usernum%
goto :eof
