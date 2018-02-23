import React from 'react';
import { Dimmer, Loader, Table } from 'semantic-ui-react';

function Tracklist(props) {
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

export default Tracklist;