<div style="text-align: center">
	<p style="font-size: 40px"> <b>Last Player On The Left </b></p>
	<img alt="GitHub" src="https://img.shields.io/github/license/Wombatlord/last-player-on-the-left?logo=Github&logoColor=green">
	<img alt="GitHub" src="https://img.shields.io/github/last-commit/Wombatlord/last-player-on-the-left?color=purple&logo=github&logoColor=purple">
	<img alt="GitHub" src="https://img.shields.io/github/languages/top/Wombatlord/last-player-on-the-left?label=Go&logo=go">
	<img alt="GitHub" src="https://img.shields.io/github/go-mod/go-version/Wombatlord/last-player-on-the-left?logo=go">
</div>

A lightweight & terminal based podcast player written in pure Go.
No external media player (vlc, mpv) required!

- Terminal UI
- Subscribe to RSS feeds
- Stream episodes from subscribed feeds.

## Installation
To build from source, clone the repo and run `go build` in the project root.

## Usage

### Subscribing to a feed
Subscribing to an RSS Feed saves the provided feed alias & url in **config.yaml**.
Aliases & Feeds can be manually added, or passed through on the command line.

Adding a feed via the commandline will then run Last Player, if you want to add multiple feeds & aliases at once, this should currently be done manually in **config.yaml**.

Alias is a user provided name for a given podcast.

`./last-player-on-the-left.exe $ALIAS -s $URL`

The following will subscribe to Last Podcast On The Left, associating the alias LPOTL to the feed and saving this information in config.yaml.

`./last-player-on-the-left.exe LPOTL -s https://feeds.simplecast.com/dCXMIpJz`

### UI & Playback Controls
Once Last Player is running, key presses will be passed through to the panel with focus.

- `TAB` or `Left / Right Arrow Keys` will change focus between panels.
- `Enter` will interact with a highlighted element in a panel:
	- If the `Podcasts` panel has focus, `Enter` will populate the `Episodes` panel.
	- If the `Episodes` panel has focus, `Enter` will begin playback of the selected episode.
- `P` will pause or resume the currently playing episode, doesn't depend on focus.