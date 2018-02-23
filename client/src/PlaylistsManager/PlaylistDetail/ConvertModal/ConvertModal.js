import React, { Component } from 'react';
import {  Button, Modal } from 'semantic-ui-react';

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

export default ConvertModal;