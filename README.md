# Playlist Exchange

Playlist Exchange (originally SpotiTube) is a web application that converts Spotify playlists into YouTube music video playlists. The site currently lives on [Heroku](https://pacific-ravine-30937.herokuapp.com/).

The API server is written in Go, using [`gorilla/mux`](https://github.com/gorilla/mux) as the HTTP router and Redis for session handling. The client is a React single-page application written with TypeScript.

## Local development environment setup

This guide assumes that you have [Redis](https://redis.io/topics/quickstart) installed locally and your [Go development environment](https://golang.org/doc/install) properly setup. Running the app also requires both Spotify and YouTube API access.

### Source code

To retrieve the repository, run:

```
go get github.com/cschen13/spotitube
```

This project manages dependencies with [Godep](https://github.com/tools/godep). Refer to the documentation there if you decide to add dependencies. Someday, we should probably transition to a different Go dependency manager, but today is not that day.

### Spotify credentials

1. Sign in to the [Spotify developers dashboard](https://developer.spotify.com/dashboard) with your Spotify account.
2. Create a new app and add the Client ID and Client Secret to your environment with variable names SPOTIFY_ID and SPOTIFY_SECRET, respectively.
3. Add `http://localhost:8081/callback/spotify` to the list of Redirect URIs for your Spotify app.

### YouTube credentials

1. Sign in to the [Google Developers Console](https://console.developers.google.com/) with your Google account.
2. Create a new project and enable the YouTube Data API v3.
3. Navigate to the Credentials page and create a new OAuth client ID.
4. Add `http://localhost:8081/callback/youtube` to the list of Authorized redirect URIs.
5. Download the JSON secret file, then save the value as an environment variable named YOUTUBE_SECRET.

### Running the app

First, get Redis up and running. From the project root:

```
$ redis-server
```

Make sure the server is running on port 6379. Now, open another terminal instance at the directory root and run the Go server:

```
$ go build
$ ./spotitube
```

Finally, open a third terminal at the directory root, then run the React app:

```
$ cd client
$ npm install
$ npm start
```

Since the Webpack development server proxies API requests to the local instance of our server, hot reloading capabilities are preserved for front-end changes. The application should now be live at `http://localhost:3000`.

### Server live reloading

To allow live reloading of our server when we make changes, install [codegangsta/gin](https://github.com/codegangsta/gin):

```
go get github.com/codegangsta/gin
```

Then, when it comes time to run the Go server, instead of running `go build` and whatnot, you can run:

```
gin -p 8080 -a 8081 ./spotitube
```

Changes to the server will now be automatically detected and rebuilt on the fly.

## End goal

Eventually, the app should ideally allow users to convert playlists in both directions and across more than just the YouTube and Spotify ecosystems (Apple Music 👀).
