With CreateObject("WScript.Shell") 
  .Exec("goland64 C:\Users\a_s_b\source\raft-db") 
  wsh.Sleep 20000 : .SendKeys "+{F10}" 
  wsh.Quit 1 
End With 
