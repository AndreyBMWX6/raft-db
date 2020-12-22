@echo off
taskkill /F /FI "imagename eq ___*" /IM ___*
taskkill /F /IM goland64.exe
exit