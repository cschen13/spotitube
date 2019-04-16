import React from "react";
import { Dimmer, Loader, Table } from "semantic-ui-react";
import { ITrack } from "../../../../../services/TrackService";

interface ITrackListProps {
  tracks: ITrack[];
}

const Tracklist: React.FunctionComponent<ITrackListProps> = ({ tracks }) => {
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
        {tracks.map((track, idx) => (
          <Table.Row key={idx}>
            <Table.Cell collapsing>{idx + 1}</Table.Cell>
            <Table.Cell collapsing>{track.title}</Table.Cell>
            <Table.Cell>{track.artist}</Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table>
  );
};

export default Tracklist;
