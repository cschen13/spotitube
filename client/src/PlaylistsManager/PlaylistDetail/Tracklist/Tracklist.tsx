import React from "react";
import { Dimmer, Loader, Table } from "semantic-ui-react";
import { ITrack } from "../../../services/TrackService";

interface ITrackListProps {
  tracks: ITrack[];
}

const Tracklist: React.FunctionComponent<ITrackListProps> = ({ tracks }) => {
  if (tracks.length === 0) {
    return (
      <Dimmer active inverted>
        <Loader
          inverted
          content="Loading Playlist (Long tracklists may take a while...)"
        />
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
        {tracks.map((track, index) => (
          <Table.Row>
            <Table.Cell collapsing>{index + 1}</Table.Cell>
            <Table.Cell collapsing>{track.title}</Table.Cell>
            <Table.Cell>{track.artist}</Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table>
  );
};

export default Tracklist;
