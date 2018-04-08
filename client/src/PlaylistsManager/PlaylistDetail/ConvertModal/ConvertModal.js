import React, { Component } from 'react';
import {  Button, Modal, Progress } from 'semantic-ui-react';

class ConvertModal extends Component {
  render() {
    const loggedInYouTube = this.props.loggedInYouTube;
    const percentProgress = this.props.percentProgress;
    const convertFailures = this.props.convertFailures;

    const content = loggedInYouTube ? (
      <div>
        <Progress
          percent={percentProgress}
          indicating={percentProgress < 100}
          success={percentProgress === 100 && convertFailures.length === 0}
          error={percentProgress === 100 && convertFailures.length > 0}>
          {
            (convertFailures.length > 0)
            ? <p>Failed to convert {this.props.convertFailures.length} tracks.</p>
            : null
          }
        </Progress>
        <p>Converting {this.props.currentTrack}...</p>
      </div>
    ) : (
      <p>To begin conversion, you must first <a href={(process.env.SPOTITUBE_HOST ? '' : 'http://localhost:8080') + '/login/youtube?returnURL=' + encodeURIComponent(window.location.pathname + window.location.search)}>login with Google/YouTube</a>.</p>
    );

    return (
      <Modal
        trigger={<Button primary onClick={this.props.onClick}>Convert to YouTube</Button>}>
        <Modal.Header>Convert Playlist to YouTube</Modal.Header>
        <Modal.Content>
          {content}
        </Modal.Content>
      </Modal>
    )
  }
}

export default ConvertModal;