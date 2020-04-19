# Musync

**Musync** is an application to sync music files using `youtube-dl`. **Musync** is
able to automatically split albums into their songs using `ffmpeg`

## Installing

The easiest way to install is to use `go get`:

```
go get github.com/cronokirby/musync
```

If you've installed the `go` programming language, this will download, compile,
and then place this program in an accessible location in your `PATH`.

In addition to this, you need to have `youtube-dl` and `ffmpeg` accessible from the command
line, since **Musync** will call each of these programs.

## Config file

**Musync** will take a description of the sources you'd like to download, and will then
automatically download these sources, and then split them into the songs they're made of.

This information is structured as a TOML file:

```toml
[[source]]
  name = "1982"
  artist = "Haircuts for Men"
  url = "https://www.youtube.com/watch?v=HSp0E0kCzVc"
  path = "vapor/chirpy/"
  timestamps = ["00:00","04:18","08:47","13:15","17:25"]
  namestamps = ["Henry's Lunch Money","Car Key Jingle","Weakling Heart","Acceptance","Midnight Luxxury"]
```

Each source is given using the "array of tables" syntax in TOML.
The *url* field gives us a location that `youtube-dl` can download the source from.
The *timestamps* and *namestamps* fields tell us how to split this source into individual songs.
The *path* field tells us where to place this source after downloading and splitting it.
This path is relative to where musync is run.
The *name* and *artist* fields are just metadata.

In this case, running `musync sync example.toml` will create a directory `vapor/chirpy/1982/` relative
to where the command is run. This directory will contain `Henry's Lunch Money.mp3` along with
the other songs. The program will also fill in the metadata, such as the album name (1982)
and the artist.

The file `example.toml` contains an example of what this file looks like.

## Adding new entries

In addition to entering new entries into the file manually, you can also use `musync`
to add them through an interactive prompt. Try running the `musync add file.toml` command.

# Usage

```
Usage: musync <command>

Flags:
  --help    Show context-sensitive help.

Commands:
  sync <path>
    Sync music files

  add <path>
    Add a new album or song to the syncing file

Run "musync <command> --help" for more information on a command.
```
