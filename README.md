# Last Player On The Left
A lightweight & terminal based podcast player written in pure Go.
No external media player (vlc, mpv) required!

- Subscribe to RSS feeds
- Stream episodes from subscribed feeds.
- Play local audio files.

## Installation
To build from source, clone the repo and run `go build` in the project root.

## Usage

### Subscribing to a feed
Subscribing to an RSS Feed saves the provided feed alias & url in config.yaml.
Aliases & Feeds can be manually added, or passed through on the command line.

Aliases can then be used from the command line in combination with flags, eg. `-l`.

Alias is a user provided name for a given podcast.

`./last-player-on-the-left.exe $ALIAS -s $URL`

The following will subscribe to Last Podcast On The Left, associating the alias LPOTL to the feed and saving this information in config.yaml.

`./last-player-on-the-left.exe LPOTL -s https://feeds.simplecast.com/dCXMIpJz`

### Play the latest episode
Grabs the feed associated to the given alias from the subscription list and begins playing the latest episode in the feed, direct from the command line.

`./last-player-on-the-left.exe $ALIAS -l`

### Play a specific episode
Episodes are indexed from the latest episode to the first episode in a feed. Pass a number greater than 0 to play a specific episode, working back from the latest.

`./last-player-on-the-left.exe LPOTL -e 1` plays the episode before the latest episode.
