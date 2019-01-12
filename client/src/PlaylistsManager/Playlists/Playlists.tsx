import React from 'react';
import { Header, List } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

interface IPlaylistsProps {
  playlists: Array<{
    ID: string;
    OwnerID: string;
    Name: string;
  }>;
}

const Playlists: React.FunctionComponent<IPlaylistsProps> = (props) => {
  return (
    <div>
      <Header as="h2">Select a playlist</Header>
      <List celled>
        {props.playlists.map((playlist) =>
          <List.Item key={playlist.ID}>
            <List.Content>
              <Link to={{pathname: `${playlist.OwnerID}/${playlist.ID}`}}>{playlist.Name}</Link>
            </List.Content>
          </List.Item>
          )}
      </List>
    </div>
  );
}

export default Playlists;