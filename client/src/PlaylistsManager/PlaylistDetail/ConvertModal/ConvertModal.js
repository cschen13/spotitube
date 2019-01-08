import React, { Component } from 'react';
import {  Button, Modal, Progress } from 'semantic-ui-react';

class ConvertModal extends Component {
  constructor(props) {
    super(props);
    this.state = {
      converted: false,
      convertFailures: [],
      currentTrack: '',
      loggedInYouTube: true,
      percentProgress: 0,
      open: false,
      playlistUrl: undefined,
    };
  }

  handleConvertClick() {
    const open = this.state.open;
    if (!open) {
      this.convertTracks()
      .then((res) => {
        this.setState({ converted: true });
        res.json().then((playlistInfo) => {
            this.setState({ playlistUrl: playlistInfo.URL });
        });
      })
      .catch((err) => {
        // Catch YouTube login errors
        console.error(err);
      });
    }
    
    this.setState({open: true});
  }

  convertTracks() {
    const ownerId = this.props.ownerId;
    const playlistId = this.props.playlistId;
    const tracks = this.props.tracks;

    let counter = 0;
    return tracks.reduce((promise, track) => {
      return promise
        .then((res) => { if (res) return res.json(); })
        .then((newPlaylist) => {
          console.log(`Converting ${track.Title}`);
          this.setState({ currentTrack: track.Title });
          let newPlaylistQuery = '';
          if (newPlaylist) {
            newPlaylistQuery = `?newPlaylistId=${newPlaylist.ID}`;
          }

          return fetch(`/playlists/${ownerId}/${playlistId}/tracks/${track.ID}${newPlaylistQuery}`, {
            credentials: 'include',
            headers: { 'Accept': 'application/json' },
            method: 'POST',
          })
          .then((res) => {
            this.handleConvertErrors(res, track);
            this.setState({ percentProgress: ++counter / tracks.length * 100 });
            return res;
          });
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

    return response;
  }

  render() {
    const loggedInYouTube = this.state.loggedInYouTube;
    const percentProgress = this.state.percentProgress;
    const convertFailures = this.state.convertFailures;
    const converted = this.state.converted;
    const playlistUrl = this.state.playlistUrl;

    const content = loggedInYouTube ? (
      <div>
        <Progress
          percent={percentProgress}
          indicating={percentProgress < 100}
          success={percentProgress === 100 && convertFailures.length === 0}
          error={percentProgress === 100 && convertFailures.length > 0}>
          {
            (convertFailures.length > 0)
            ? <p>Failed to convert {this.state.convertFailures.length} tracks.</p>
            : null
          }
        </Progress>
        
        {
            converted
            ? <p>Your playlist has been converted! See it <a href={playlistUrl}>here</a>.</p>
            : <p>Converting {this.state.currentTrack}...</p>
        }
      </div>
    ) : (
      <p>To begin conversion, you must first <a href={(process.env.REACT_APP_SPOTITUBE_HOST ? '' : 'http://localhost:8080') + '/login/youtube?returnURL=' + encodeURIComponent(window.location.pathname + window.location.search)}>login with Google/YouTube</a>.</p>
    );

    return (
      <Modal
        closeOnRootNodeClick={false}
        trigger={<Button primary>Convert to YouTube</Button>}
        onOpen={() => this.handleConvertClick()}>
        <Modal.Header>Convert Playlist to YouTube</Modal.Header>
        <Modal.Content>
          {content}
        </Modal.Content>
      </Modal>
    );
  }
}

export default ConvertModal;