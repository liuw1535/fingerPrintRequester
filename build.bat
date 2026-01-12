@echo off
setlocal

echo Building TLS Fingerprint Requester...

REM Download dependencies
echo Running go mod tidy...
go mod tidy

REM Create bin directory
if not exist bin mkdir bin

REM Build flags for size optimization
set LDFLAGS=-s -w

REM Linux AMD64
echo Building for Linux AMD64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="%LDFLAGS%" -o bin/fingerprint_linux_amd64 ./cmd/tlsRequester
if errorlevel 1 goto error

REM Windows AMD64
echo Building for Windows AMD64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="%LDFLAGS%" -o bin/fingerprint_windows_amd64.exe ./cmd/tlsRequester
if errorlevel 1 goto error

REM Android ARM64
echo Building for Android ARM64...
set GOOS=android
set GOARCH=arm64
go build -ldflags="%LDFLAGS%" -o bin/fingerprint_android_arm64 ./cmd/tlsRequester
if errorlevel 1 goto error

echo Build completed successfully!
echo Binaries are in the bin/ directory
goto end

:error
echo Build failed!
exit /b 1

:end
endlocal
