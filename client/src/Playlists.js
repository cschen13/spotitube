import React, { Component } from 'react';
import { Header, List, Image, Button, Modal, Table, Loader, Dimmer } from 'semantic-ui-react';
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
          <List.Item key={playlist.ID}>
            <List.Content>
              <Link to={{pathname: playlist.OwnerID + '/' + playlist.ID}}>{playlist.Name}</Link>
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
      hasError: false,
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
      return res.json();
    })
    .then((playlist) => {
      console.log(playlist);
      this.setState({
        playlist: {
          name: playlist.Name,
          url: playlist.URL,
          coverUrl: playlist.CoverURL,
        },
        showConvertModal: false,
      });

      return fetch(`/playlists/${ownerId}/${playlistId}/tracks`, {
        credentials: 'include',
        headers: headers
      });
    })
    .then((res) => {
      console.log(res);
      return res.json();
    })
    .then((tracks) => {
      console.log(tracks);
      this.setState({
        tracks: tracks,
      });
    })
    .catch((err) => {
      this.setState({
        hasError: true,
      });
      console.log('Request failed', err);
    })
  }

  handleConvertClick() {
    this.setState({showConvertModal: true});
    console.log("Convert button clicked");
  }

  render() {
    const playlistName = this.state.playlist.name;
    const coverUrl = this.state.playlist.coverUrl;
    const tracks = this.state.tracks;

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        {this.state.hasError
          ? <p>An error occurred retrieving some part of the playlist.</p>
          : <div>
              <Image src={coverUrl === '' ? noArtwork : coverUrl} size="medium" />
              <ConvertModal onClick={() => this.handleConvertClick()} />
              <TrackList tracks={tracks} />
            </div>
        }
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

function TrackList(props) {
  if (props.tracks.length === 0) {
    return (
      <Dimmer active inverted>
        <Loader inverted content='Loading Playlist (Long tracklists may take a while...)' />
      </Dimmer>
    );
  }

  return (
    <Table definition>
      <Table.Header>
        <Table.Row>
          <Table.HeaderCell />
          <Table.HeaderCell>Title</Table.HeaderCell>
          <Table.HeaderCell>Artist</Table.HeaderCell>
        </Table.Row>
      </Table.Header>

      <Table.Body>
        {props.tracks.map((track, index) => 
          <Table.Row>
            <Table.Cell collapsing>{index+1}</Table.Cell>
            <Table.Cell collapsing>{track.Title}</Table.Cell>
            <Table.Cell>{track.Artist}</Table.Cell>
          </Table.Row>
        )}
      </Table.Body>
    </Table>
  );
}

export default PlaylistsManager;