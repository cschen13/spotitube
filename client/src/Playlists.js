import React, { Component } from 'react';
import { Header, List, Image, Button } from 'semantic-ui-react';
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import noArtwork from './no-artwork.png';

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

function Playlists(props) {
  return (
    <div>
      <Header as="h2">Select a playlist</Header>
      <List celled>
        {props.playlists.map(playlist => 
          <List.Item key={playlist.id}>
            <List.Content>
              <Link to={{pathname: playlist.ownerId + '/' + playlist.id}}>{playlist.name}</Link>
            </List.Content>
          </List.Item>
          )}
      </List>
    </div>
  );
}

class PlaylistDetail extends Component {
  constructor(props) {
    super(props);
    this.state = {
      playlist: {
        name: '',
        url: '#',
        coverUrl: '',
      },
      tracks: [],
    };
  }

  componentDidMount() {
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;

    fetch('/playlists/' + ownerId + '/' + playlistId)
    .then((res) => {
      console.log(res);
      if (res.ok) {
        res.json().then((json) => {
          this.setState({
              playlist: {
                name: json.playlist.name,
                url: json.playlist.url,
                coverUrl: json.playlist.coverUrl,
            },
            tracks: json.tracks,
          });
        });
      } else {
        // TODO: Handle errors
      }
    });
  }

  render() {
    const playlistName = this.state.playlist.name;
    const coverUrl = this.state.playlist.coverUrl;

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        <Image src={coverUrl === '' ? noArtwork : coverUrl} size="medium" />
        <p>TODO: Track Listing Goes Here</p>
        <Button primary>Convert to YouTube</Button>
      </div>
    );
  }
}

export default PlaylistsManager;