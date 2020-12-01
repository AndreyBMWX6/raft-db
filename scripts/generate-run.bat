@echo off
rem making run.vbs scripts for each server

echo With CreateObject("WScript.Shell") > C:\Users\a_s_b\source\raft-db\scripts\run.vbs
echo   .Exec("goland64 C:\Users\a_s_b\source\raft-db") >> C:\Users\a_s_b\source\raft-db\scripts\run.vbs
echo   wsh.Sleep 20000 : .SendKeys "+{F10}" >> C:\Users\a_s_b\source\raft-db\scripts\run.vbs
echo   wsh.Quit 1 >> C:\Users\a_s_b\source\raft-db\scripts\run.vbs
echo End With >> C:\Users\a_s_b\source\raft-db\scripts\run.vbs

echo With CreateObject("WScript.Shell") > C:\Users\a_s_b\source\"raft-db node 2"\scripts\run.vbs
echo   .Exec("goland64 C:\Users\a_s_b\source\raft-db node 2") >> C:\Users\a_s_b\source\"raft-db node 2"\scripts\run.vbs
echo   wsh.Sleep 20000 : .SendKeys "+{F10}" >> C:\Users\a_s_b\source\"raft-db node 2"\scripts\run.vbs
echo   wsh.Quit 1 >> C:\Users\a_s_b\source\"raft-db node 2"\scripts\run.vbs
echo End With >> C:\Users\a_s_b\source\"raft-db node 2"\scripts\run.vbs

echo With CreateObject("WScript.Shell") > C:\Users\a_s_b\source\"raft-db node 3"\scripts\run.vbs
echo   .Exec("goland64 C:\Users\a_s_b\source\raft-db node 3") >> C:\Users\a_s_b\source\"raft-db node 3"\scripts\run.vbs
echo   wsh.Sleep 20000 : .SendKeys "+{F10}" >> C:\Users\a_s_b\source\"raft-db node 3"\scripts\run.vbs
echo   wsh.Quit 1 >> C:\Users\a_s_b\source\"raft-db node 3"\scripts\run.vbs
echo End With >> C:\Users\a_s_b\source\"raft-db node 3"\scripts\run.vbs

echo With CreateObject("WScript.Shell") > C:\Users\a_s_b\source\"raft-db node 4"\scripts\run.vbs
echo   .Exec("goland64 C:\Users\a_s_b\source\raft-db node 4") >> C:\Users\a_s_b\source\"raft-db node 4"\scripts\run.vbs
echo   wsh.Sleep 20000 : .SendKeys "+{F10}" >> C:\Users\a_s_b\source\"raft-db node 4"\scripts\run.vbs
echo   wsh.Quit 1 >> C:\Users\a_s_b\source\"raft-db node 4"\scripts\run.vbs
echo End With >> C:\Users\a_s_b\source\"raft-db node 4"\scripts\run.vbs

echo With CreateObject("WScript.Shell") > C:\Users\a_s_b\source\"raft-db node 5"\scripts\run.vbs
echo   .Exec("goland64 C:\Users\a_s_b\source\raft-db node 5") >> C:\Users\a_s_b\source\"raft-db node 5"\scripts\run.vbs
echo   wsh.Sleep 20000 : .SendKeys "+{F10}" >> C:\Users\a_s_b\source\"raft-db node 5"\scripts\run.vbs
echo   wsh.Quit 1 >> C:\Users\a_s_b\source\"raft-db node 5"\scripts\run.vbs
echo End With >> C:\Users\a_s_b\source\"raft-db node 5"\scripts\run.vbs

echo With CreateObject("WScript.Shell") > C:\Users\a_s_b\source\"raft-db node 6"\scripts\run.vbs
echo   .Exec("goland64 C:\Users\a_s_b\source\raft-db node 6") >> C:\Users\a_s_b\source\"raft-db node 6"\scripts\run.vbs
echo   wsh.Sleep 20000 : .SendKeys "+{F10}" >> C:\Users\a_s_b\source\"raft-db node 6"\scripts\run.vbs
echo   wsh.Quit 1 >> C:\Users\a_s_b\source\"raft-db node 6"\scripts\run.vbs
echo End With >> C:\Users\a_s_b\source\"raft-db node 6"\scripts\run.vbs

exit