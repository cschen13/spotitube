import * as React from "react";
import { IPlaylist } from "../../services/PlaylistService";
import { ITrack } from "../../services/TrackService";
import ConvertPanel from "./ConvertPanel/ConvertPanel";
import PlaylistsManager from "./PlaylistsManager/PlaylistsManager";

export interface ILoadedPlaylist {
  details?: IPlaylist;
  tracks?: ITrack[];
}

interface IConverterProps {
  hasGetError: boolean;
  playlists: IPlaylist[];
}

interface IConverterState {
  readonly playlistsToConvert: ILoadedPlaylist[];
}

export default class Converter extends React.Component<
  IConverterProps,
  IConverterState
> {
  constructor(props: IConverterProps) {
    super(props);
    this.state = {
      playlistsToConvert: []
    };
  }

  public render() {
    const { hasGetError, playlists } = this.props;
    const { playlistsToConvert } = this.state;
    return (
      <div>
        <PlaylistsManager
          playlists={playlists}
          hasGetError={hasGetError}
          onConvert={(playlist, tracks) => this.handleConvert(playlist, tracks)}
        />
        <ConvertPanel playlistsToConvert={playlistsToConvert} />
      </div>
    );
  }

  private handleConvert(playlist?: IPlaylist, tracks?: ITrack[]) {
    const { playlistsToConvert } = this.state;
    this.setState({
      playlistsToConvert: [...playlistsToConvert, { details: playlist, tracks }]
    });
  }
}
