import React, { Component } from 'react';
import { Header, Image } from 'semantic-ui-react';
import noArtwork from '../../imgs/no-artwork.png';
import ConvertModal from './ConvertModal/ConvertModal';
import Tracklist from './Tracklist/Tracklist';

class PlaylistDetail extends Component {
  constructor(props) {
    super(props);
    this.state = {
      converted: false,
      convertFailures: [],
      hasGetError: false,
      loggedInYouTube: true,
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
    const showConvertModal = this.state.showConvertModal;
    if (!showConvertModal) {
      this.convertTracks()
      .then((res) => {
        this.setState({ converted: true });
      })
      .catch((err) => {
        // Catch YouTube login errors
        console.error(err);
      })
    }
    
    this.setState({showConvertModal: true});
  }

  convertTracks() {
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;
    const tracks = this.state.tracks;

    return tracks.reduce((promise, track) => {
      return promise
        .then((res) => {
          console.log(`Converting ${track.Title}`);
          return fetch(`/playlists/${ownerId}/${playlistId}/tracks/${track.ID}`, {
            credentials: 'include',
            headers: {
              'Accept': 'application/json',
            },
            method: 'POST',
          })
          .then((res) => {
            if (!res.ok) {
              if (res.status === 401) {
                this.setState({ loggedInYouTube: false });
                throw new Error(res.statusText);
              } else {
                console.log(`Failed to convert ${track.Title}`);
                this.setState({ convertFailures: [...this.state.convertFailures, track] });
              }
            }
          });
        });
    }, Promise.resolve());
  }

  render() {
    const convertFailures = this.state.convertFailures;
    const loggedInYouTube = this.state.loggedInYouTube;
    const playlistName = this.state.playlist.name;
    const coverUrl = this.state.playlist.coverUrl;
    const tracks = this.state.tracks;

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        {this.state.hasGetError
          ? <p>An error occurred retrieving some part of the playlist.</p>
          : <div>
              <Image src={coverUrl === '' ? noArtwork : coverUrl} size="medium" />
              <ConvertModal
                loggedInYouTube={loggedInYouTube} 
                convertFailures={convertFailures}
                onClick={() => this.handleConvertClick()} />
              <Tracklist tracks={tracks} />
            </div>
        }
      </div>
    );
  }
}

export default PlaylistDetail;