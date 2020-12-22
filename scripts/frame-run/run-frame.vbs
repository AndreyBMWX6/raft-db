With CreateObject("WScript.Shell") 
  .Exec("goland64 ..\..\raft-db") 
  wsh.Sleep 20000 : .SendKeys "+{F10}" 
  wsh.Quit 1 
End With 
