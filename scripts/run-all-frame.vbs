With CreateObject("WScript.Shell") 
  .CurrentDirectory="..\..\raft-db"
  .Exec("goland64 .")
  wsh.Sleep 20000 : .SendKeys "+{F10}"
End With

With CreateObject("WScript.Shell") 
  .CurrentDirectory="..\raft-db node 2"
  .Exec("goland64 .")
  wsh.Sleep 9000 : .SendKeys "+{F10}"
End With

With CreateObject("WScript.Shell") 
  .CurrentDirectory="..\raft-db node 3"
  .Exec("goland64 .")
  wsh.Sleep 9000 : .SendKeys "+{F10}"
End With

With CreateObject("WScript.Shell") 
  .CurrentDirectory="..\raft-db node 4"
  .Exec("goland64 .")
  wsh.Sleep 9000 : .SendKeys "+{F10}"
End With

With CreateObject("WScript.Shell") 
  .CurrentDirectory="..\raft-db node 5"
  .Exec("goland64 .")
  wsh.Sleep 9000 : .SendKeys "+{F10}"
End With

With CreateObject("WScript.Shell") 
  .CurrentDirectory="..\raft-db node 6"
  .Exec("goland64 .")
  wsh.Sleep 9000 : .SendKeys "+{F10}"
  wsh.Quit 1
End With