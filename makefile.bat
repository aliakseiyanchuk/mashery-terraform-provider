@echo off

set PLUGIN_DIR=%APPDATA%\terraform.d\plugins\yanchuk.nl\aliakseiyanchuk\mashery\0.4\windows_amd64
mkdir %PLUGIN_DIR%
set WIN_BINARY_NAME=terraform-provider-mashery.exe

go build -o %WIN_BINARY_NAME%
copy %WIN_BINARY_NAME% %PLUGIN_DIR%
