import React, { Component } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { IPlaylist } from "../services/PlaylistService";
import PlaylistDetail from "./PlaylistDetail/PlaylistDetail";
import Playlists from "./Playlists/Playlists";

interface IPlaylistsManagerProps {
  playlists: IPlaylist[];
}

class PlaylistsManager extends React.Component<IPlaylistsManagerProps> {
  public render() {
    return (
      <Router>
        <div>
          <Switch>
            <Route
              exact
              path="/"
              render={() => <Playlists playlists={this.props.playlists} />}
            />
            <Route path="/:ownerId/:playlistId" component={PlaylistDetail} />
          </Switch>
        </div>
      </Router>
    );
  }
}

export default PlaylistsManager;
