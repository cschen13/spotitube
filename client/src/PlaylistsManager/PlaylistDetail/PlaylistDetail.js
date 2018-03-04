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
    .then((res) => res.json())
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
    .then((res) => res.json())
    .then((tracks) => { this.setState({ tracks: tracks }); })
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
        .then((res) => { if (res) return res.json(); })
        .then((newPlaylist) => {
          console.log(`Converting ${track.Title}`);
          let newPlaylistQuery = '';
          if (newPlaylist) {
            newPlaylistQuery = `?newPlaylistId=${newPlaylist.ID}`;
          }

          return fetch(`/playlists/${ownerId}/${playlistId}/tracks/${track.ID}${newPlaylistQuery}`, {
            credentials: 'include',
            headers: { 'Accept': 'application/json' },
            method: 'POST',
          })
          .then((res) => this.handleConvertErrors(res, track));
        });
    }, Promise.resolve());
  }

  handleConvertErrors(response, track) {
    if (!response.ok) {
      if (response.status === 401) {
        this.setState({ loggedInYouTube: false });
        throw new Error(response.statusText);
      } else {
        console.error(`Failed to convert ${track.Title}`);
        this.setState({ convertFailures: [...this.state.convertFailures, track] });
      }
    }
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