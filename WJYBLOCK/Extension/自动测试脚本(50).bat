@echo off

echo  ?????????

cd /d %~dp0

echo ????????????

del /a %USERPROFILE%\BlockChain

echo ?????????????

xcopy BlockChain %USERPROFILE%\BlockChain

set times=50

start server.exe
TIMEOUT /T 2

for /l %%i in (1,1,%times%) do (

start client.exe -v false
TIMEOUT /T 1

)
