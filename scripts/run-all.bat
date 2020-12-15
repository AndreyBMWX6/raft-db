@echo off

set ip=%~1
set firstIpPort=%~2
set lastIpPort=%~3
set firstUrlPort=%~4
set lastUrlPort=%~5
set /a usercnt=(%lastIpPort% - %firstIpPort%) + 1
set /a usernum=1
set "router=true"
set "norouter=false"

echo running %usercnt% servers on ip:%ip%
echo ip=%ip%
echo firstIpPort=%firstIpPort%
echo lastIpPort=%lastIpPort%
echo firstUrlPort=%firstUrlPort%
echo lastUrlPort=%lastUrlPort%

for /l %%i in (1, 1, %usercnt%) do (
  if "%%i" == "1" (
  call :runtrue
 )

 if not "%%i" == "1" (
  call :runfalse
 )
call :inc
)
goto :eof

:inc
rem echo fisrtIpPort=%firstIpPort%
set /a firstIpPort=firstIpPort + 1
rem echo firstIpPort=%firstIpPort%
rem echo firstURLPort=%firstURLPort%
set /a firstURLPort=firstURLPort + 1
rem echo firstURLPort=%firstURLPort%
rem echo usernum=%usernum%
set /a usernum=usernum + 1
rem echo usernum=%usernum%
goto :eof

:runtrue
rem start go run ../../raft-db/cmd/main.go -all %ip% %firstIpPort% %firstUrlPort% %router% user%usernum%
start run.bat %ip% %firstIpPort% %firstUrlPort% %router% user%usernum%
goto :eof

:runfalse
rem start go run ../../raft-db/cmd/main.go -all %ip% %firstIpPort% %firstUrlPort% %norouter% user%usernum%
start run.bat %ip% %firstIpPort% %firstUrlPort% %norouter% user%usernum%
goto :eof