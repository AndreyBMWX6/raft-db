@echo off
cd "../../raft-db/cmd"
start go run main.go -all 127.0.0.1 8001 8081
exit