import React, { Component } from 'react';
import './App.css';
import { Header } from 'semantic-ui-react'

class App extends Component {

  render() {
    console.log(process.env.SPOTITUBE_HOST);

    return (
      <div className="ui main container">
        <header>
          <Header as="h1">SpotiTube</Header>
        </header>
        <p>
          Convert your Spotify playlists to YouTube music video playlists.
        </p>
        <p>
          <a href={(process.env.SPOTITUBE_HOST ? '' : 'http://localhost:8080') + '/login/spotify'}>
            Login with Spotify
          </a> to get started.
        </p>
      </div>
    );
  }
}

export default App;
