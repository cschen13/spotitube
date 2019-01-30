import React from "react";
import { Link } from "react-router-dom";
import { Header, List } from "semantic-ui-react";
import { IPlaylist } from "../../services/PlaylistService";

interface IPlaylistsProps {
  playlists: IPlaylist[];
}

const Playlists: React.FunctionComponent<IPlaylistsProps> = props => {
  return (
    <div>
      <Header as="h2">Select a playlist</Header>
      <List celled>
        {props.playlists.map(playlist => (
          <List.Item key={playlist.id}>
            <List.Content>
              <Link to={{ pathname: `${playlist.ownerId}/${playlist.id}` }}>
                {playlist.name}
              </Link>
            </List.Content>
          </List.Item>
        ))}
      </List>
    </div>
  );
};

export default Playlists;
