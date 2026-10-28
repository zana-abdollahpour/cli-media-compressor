# imgcompress

`imgcompress` is a command-line tool that scans the current directory, compresses supported image and video files, and writes the results under an output directory while preserving their relative paths. Original files are left in place. If compression would make a file larger, the original bytes are copied to the output instead.

## Supported formats

- Images: JPEG (`.jpg`, `.jpeg`), PNG, and GIF
- Videos: MP4, M4V, MOV, AVI, MKV, and WebM

Image quality applies to JPEG encoding. PNG and GIF are re-encoded with their respective encoders. Video compression uses H.264/AAC for most formats and VP9/Opus for WebM.

## Build

Requirements:

- Go 1.26.5 or newer
- For video compression, FFmpeg on `PATH`, unless you build the Windows x64 executable with the bundled FFmpeg option below

Build a regular executable:

```sh
go build -o imgcompress .
```

On Windows, use `go build -o imgcompress.exe .`. This normal build does not include FFmpeg; video compression requires a separately installed FFmpeg available on `PATH`. Image compression works without FFmpeg.

### Windows x64 build with embedded FFmpeg

Run the included script from the repository root in a Windows command prompt:

```bat
build-windows.cmd
```

The script creates `dist\imgcompress.exe`. It uses the pinned FFmpeg 7z archive in the repository root if present; otherwise it downloads that archive. It verifies the archive and extracted executable hashes, extracts `ffmpeg.exe`, then embeds it in the Go build. The finished app extracts FFmpeg into the user's cache directory when needed, so the target Windows x64 machine does not need FFmpeg installed.

Build prerequisites for this script:

- Go 1.26.5 or newer
- `curl.exe` and PowerShell, available on supported Windows installations
- Windows `tar.exe` with 7z extraction support
- Internet access on the first build if the pinned archive is not already present

The archive and extracted executable are ignored by Git, so a fresh clone downloads the archive when you run the script. The generated executable is also ignored. This script targets Windows x64 only; other platforms use the regular Go build and an FFmpeg installation on `PATH` for video compression.

See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) for FFmpeg attribution and redistribution licensing notes.

## Use

Run the program from the directory tree you want to process:

```sh
imgcompress
```

By default, it prompts for JPEG quality from 1 to 100 (default 40), then scans the current directory recursively and writes results to `compressed/`, keeping each file's relative path. The output directory is excluded from scanning.

Options:

```text
-out string     output directory (default "compressed")
-quality int    JPEG quality from 1 to 100; supplying this skips the prompt (default 40)
-v              print input and output sizes for each file
-h              show command-line help
```

Examples:

```sh
# Use the default output directory and answer the quality prompt
imgcompress

# Set quality and output directory without prompting
imgcompress -quality 65 -out optimized

# Include per-file size details
imgcompress -quality 65 -v
```

The output directory can be relative to the current directory or an absolute path. Existing output files with matching paths may be overwritten.
