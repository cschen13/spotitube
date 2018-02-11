import React, { Component } from 'react';
import './App.css';
import { Header } from 'semantic-ui-react';
import PlaylistsManager from './Playlists.js';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loggedIn: false,
      playlists: [],
    };
  }

  componentDidMount() {
    fetch('/playlists', {
      credentials: 'include',
    })
    .then((response) => {
      // console.log(response);
      this.setState({ loggedIn: response.ok });
      if (response.ok) {
        response.json().then((playlists) => {
          console.log(playlists);
          this.setState({ playlists: playlists });
        });
      }
    });
  }

  render() {
    let landing;
    if (this.state.loggedIn) {
      landing = <PlaylistsManager playlists={this.state.playlists} />;
    } else {
      landing = <Greeting />;
    }

    return (
      <div className="ui main container">
        <header>
          <Header as="h1">SpotiTube</Header>
        </header>
        {landing}
      </div>
    );
  }
}

function Greeting(props) {
  return (
    <div>
      <p>
        Convert your Spotify playlists to YouTube music video playlists.
      </p>
      <p>
        <a href={(process.env.SPOTITUBE_HOST ? '' : 'http://localhost:8080') + '/login/spotify?returnURL=' + encodeURIComponent(window.location.pathname + window.location.search)}>
          Login with Spotify
        </a> to get started.
      </p>
    </div>
  );
}

export default App;
