import React, { Component } from 'react';
import { Header, List, Image, Button, Modal } from 'semantic-ui-react';
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import noArtwork from './imgs/no-artwork.png';

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
      loggedInYouTube: false,
    };
  }

  componentDidMount() {
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;

    let headers = new Headers();
    headers.append('Accept', 'application/json')
    fetch(`/playlists/${ownerId}/${playlistId}`, {
      credentials: 'include',
      headers: headers
    })
    .then((res) => {
      console.log(res);
      if (res.ok) {
        res.json().then((json) => {
          console.log(json);
          this.setState({
              playlist: {
                name: json.playlist.name,
                url: json.playlist.url,
                coverUrl: json.playlist.coverUrl,
            },
            tracks: json.tracks,
            showConvertModal: false,
          });
        });
      } else {
        // TODO: Handle errors
      }
    });
  }

  handleConvertClick() {
    this.setState({showConvertModal: true});
    console.log("Convert button clicked");
  }

  render() {
    const playlistName = this.state.playlist.name;
    const coverUrl = this.state.playlist.coverUrl;

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        <Image src={coverUrl === '' ? noArtwork : coverUrl} size="medium" />
        <p>TODO: Track Listing Goes Here</p>
        <ConvertModal onClick={() => this.handleConvertClick()} />
      </div>
    );
  }
}

class ConvertModal extends Component {
  render() {
    return (
      <Modal
        trigger={<Button primary onClick={this.props.onClick}>Convert to YouTube</Button>}>
        <Modal.Header>Convert Playlist to YouTube</Modal.Header>
        <Modal.Content>
          <p>To begin conversion, you must first <a href={(process.env.SPOTITUBE_HOST ? '' : 'http://localhost:8080') + '/login/youtube?returnURL=' + encodeURIComponent(window.location.pathname + window.location.search)}>login with Google/YouTube</a>.</p>
        </Modal.Content>
      </Modal>
    )
  }
}

export default PlaylistsManager;