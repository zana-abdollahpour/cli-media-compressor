@echo off
setlocal EnableExtensions

set "ARCHIVE=ffmpeg-2026-10-25-git-46f23be8b1-essentials_build.7z"
set "ARCHIVE_URL=https://www.gyan.dev/ffmpeg/builds/packages/%ARCHIVE%"
set "ARCHIVE_SHA256=055DD68F62E2F55758D9BDF51B5480CEB0172C7D05B41E088CADCEC121CEC681"
set "FFMPEG_SHA256=362A203A8B06FB5E1D7173FA24E3FD228D2640468EE01F5A3F5D2E71593A213C"
set "ASSET=internal\videos\assets\windows_amd64\ffmpeg.exe"
set "MEMBER=ffmpeg-2026-10-25-git-46f23be8b1-essentials_build/bin/ffmpeg.exe"

if not exist "%ARCHIVE%" (
    echo Downloading pinned FFmpeg archive...
    curl.exe -fL --retry 3 -o "%ARCHIVE%" "%ARCHIVE_URL%"
    if errorlevel 1 exit /b 1
)

powershell -NoProfile -Command "if ((Get-FileHash -Algorithm SHA256 -LiteralPath '%ARCHIVE%').Hash -ne '%ARCHIVE_SHA256%') { exit 1 }"
if errorlevel 1 (
    echo FFmpeg archive SHA-256 mismatch.
    exit /b 1
)

if not exist "internal\videos\assets\windows_amd64" mkdir "internal\videos\assets\windows_amd64"
"%SystemRoot%\System32\tar.exe" -xOf "%ARCHIVE%" "%MEMBER%" > "%ASSET%"
if errorlevel 1 (
    echo Could not extract ffmpeg.exe from the archive.
    exit /b 1
)

powershell -NoProfile -Command "if ((Get-FileHash -Algorithm SHA256 -LiteralPath '%ASSET%').Hash -ne '%FFMPEG_SHA256%') { exit 1 }"
if errorlevel 1 (
    echo Extracted ffmpeg.exe SHA-256 mismatch.
    exit /b 1
)

"%ASSET%" -version 2>&1 | findstr /C:"ffmpeg version 2026-10-25-git-46f23be8b1" >nul
if errorlevel 1 (
    echo Unexpected FFmpeg version.
    exit /b 1
)
"%ASSET%" -version 2>&1 | findstr /C:"--enable-gpl" >nul
if errorlevel 1 exit /b 1
"%ASSET%" -version 2>&1 | findstr /C:"--enable-libx264" >nul
if errorlevel 1 exit /b 1

if not exist dist mkdir dist
set "GOCACHE=%TEMP%\mediacompress-gocache"
if not exist "%GOCACHE%" mkdir "%GOCACHE%"
go build -tags embedffmpeg -o dist\mediacompress.exe .
if errorlevel 1 exit /b 1

echo Created self-contained executable: dist\mediacompress.exe

