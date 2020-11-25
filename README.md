# raft-db 

#### Note:
For correct importing Go modules integration should be disabled in your IDE.

## How to run cluster:
If you're running cluster on 1 device:
Firstly, you should change paths in script.bat to your cluster nodes paths.

Paths are set this way:     
cd "<"path to your raft-db folder">/cmd"

Here's an example:  
cd "C:/Users/a_s_b/source/raft-db/cmd"

So, you have to do several raft-db folders and set the path in script.bat.

Then you have to change config.go file in your raft-db copies.
You should choose an address from neighbors slice and swap it with resolving address.
#### For Windows Users:
Win + R  
$ cd <"path to raft-db folder">  
$ script.bat
