# cli_subsonic_client

Command line application to listen to music from a server that respects the subsonic standard.

The application must be launched in server mode, then we can use the same application to control the server. This allows to manage the music in command line and to integrate this music player from [polybar](https://github.com/polybar/polybar).

## Dependencies
 - [mpv](https://mpv.io/)

## Usage
### Launch the server
```
cli_subsonic server
```
### Command to control the server
```
cli_subsonic search album <album_name>
cli_subsonic add album <album_name>
cli_subsonic add playlist <playlist_name>
cli_subsonic queued
cli_subsonic current
cli_subsonic next
cli_subsonic pause
cli_subsonic play
cli_subsonic prev
cli_subsonic random
cli_subsonic repeat <on|off>
cli_subsonic random <on|off>
cli_subsonic stop
cli_subsonic status
cli_subsonic clean
cli_subsonic favorite
```
## Credit
Part of the code comes from this [repository](https://github.com/wildeyedskies/stmp).

## Contact me
dev@guillaumepin.ch