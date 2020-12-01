@echo off
cd "C:/Users/a_s_b/source/raft-db/cmd"
start go run main.go
cd "C:/Users/a_s_b/source/raft-db node 2/cmd"
start go run main.go
cd "C:/Users/a_s_b/source/raft-db node 3/cmd"
start go run main.go
cd "C:/Users/a_s_b/source/raft-db node 4/cmd"
start go run main.go
cd "C:/Users/a_s_b/source/raft-db node 5/cmd"
start go run main.go
cd "C:/Users/a_s_b/source/raft-db node 6/cmd"
start go run main.go
exit