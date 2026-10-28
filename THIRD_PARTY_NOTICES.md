# FFmpeg redistribution note

This app embeds Gyan's Windows x64 static FFmpeg Essentials build:
`2026-10-25-git-46f23be8b1`, from
`ffmpeg-2026-10-25-git-46f23be8b1-essentials_build.7z`. The build has GPL and
version 3 support enabled and includes `libx264`, which the app uses for video
encoding. The app launches FFmpeg as a separate executable after extracting it
to the user's cache directory.

Gyan's build information is at https://www.gyan.dev/ffmpeg/builds/ and FFmpeg's
licensing guidance is at https://ffmpeg.org/legal.html. Review the FFmpeg and
included codec licenses, and provide the corresponding source and notices
required by those licenses before redistributing a built app.
