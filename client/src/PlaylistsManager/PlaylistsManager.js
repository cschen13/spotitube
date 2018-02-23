import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Playlists from './Playlists/Playlists';
import PlaylistDetail from './PlaylistDetail/PlaylistDetail';

class PlaylistsManager extends Component {
  render() {
    return (
      <Router>
        <div>
          <Route exact path="/" render={() => 
            <Playlists
              playlists={this.props.playlists}
              onClick={(p) => { this.setState({ selectedPlaylist: p }) }} />
          } />
          <Route path="/:ownerId/:playlistId" component={PlaylistDetail} />
        </div>
      </Router>
    );
  }
}

export default PlaylistsManager;