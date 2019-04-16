import * as React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { IPlaylist } from "../../../services/PlaylistService";
import PlaylistDetail from "./PlaylistDetail/PlaylistDetail";
import Playlists from "./Playlists/Playlists";
import { ITrack } from "../../../services/TrackService";

interface IPlaylistsManagerProps {
  playlists: IPlaylist[];
  hasGetError: boolean;
  onConvert: (playlist?: IPlaylist, tracks?: ITrack[]) => void;
}

class PlaylistsManager extends React.Component<IPlaylistsManagerProps> {
  public render() {
    const { playlists, hasGetError } = this.props;
    if (hasGetError) {
      return (
        <p>
          An error occurred while retrieving your playlists. Try again later.
        </p>
      );
    }

    return hasGetError ? (
      <p>An error occurred while retrieving your playlists. Try again later.</p>
    ) : (
      <Router>
        <Switch>
          <Route
            exact
            path="/"
            render={() => <Playlists playlists={playlists} />}
          />
          <Route
            path="/:ownerId/:playlistId"
            render={props => (
              <PlaylistDetail
                {...props}
                onConvert={(playlist, tracks) =>
                  this.handleConvert(playlist, tracks)
                }
              />
            )}
          />
        </Switch>
      </Router>
    );
  }

  private handleConvert(playlist?: IPlaylist, tracks?: ITrack[]) {
    this.props.onConvert(playlist, tracks);
  }
}

export default PlaylistsManager;
