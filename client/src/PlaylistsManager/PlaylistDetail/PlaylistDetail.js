import React, { Component } from 'react';
import { Header, Image } from 'semantic-ui-react';
import noArtwork from '../../imgs/no-artwork.png';
import ConvertModal from './ConvertModal/ConvertModal';
import Tracklist from './Tracklist/Tracklist';

class PlaylistDetail extends Component {
  constructor(props) {
    super(props);
    this.state = {
      hasGetError: false,
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

  render() {
    const playlistName = this.state.playlist.name;
    const coverUrl = this.state.playlist.coverUrl;
    const tracks = this.state.tracks;
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        {this.state.hasGetError
          ? <p>An error occurred retrieving some part of the playlist.</p>
          : <div>
              <Image src={coverUrl === '' ? noArtwork : coverUrl} size="medium" />
              <ConvertModal
                ownerId={ownerId}
                playlistId={playlistId}
                tracks={tracks} />
              <Tracklist tracks={tracks} />
            </div>
        }
      </div>
    );
  }
}

export default PlaylistDetail;