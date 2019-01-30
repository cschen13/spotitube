import * as React from "react";
import { RouteComponentProps } from "react-router-dom";
import { Header, Image } from "semantic-ui-react";
// @ts-ignore: https://github.com/Microsoft/TypeScript/issues/15146
import noArtwork from "../../imgs/no-artwork.png";
import playlistService, { IPlaylist } from "../../services/PlaylistService";
import { ITrack } from "../../services/TrackService";
import ConvertModal from "./ConvertModal/ConvertModal";
import Tracklist from "./Tracklist/Tracklist";

// TODO: Pass playlist details as props instead of using API to set state
interface IPlaylistDetailState {
  readonly hasGetError: boolean;
  readonly playlist?: IPlaylist;
  readonly tracks?: ITrack[];
}

interface IPlaylistDetailMatchProps {
  ownerId: string;
  playlistId: string;
}

type PlaylistDetailProps = RouteComponentProps<IPlaylistDetailMatchProps>;
class PlaylistDetail extends React.Component<
  PlaylistDetailProps,
  IPlaylistDetailState
> {
  constructor(props: PlaylistDetailProps) {
    super(props);
    this.state = {
      hasGetError: false,
      playlist: {
        coverUrl: "",
        name: ""
      } as IPlaylist,
      tracks: undefined
    };
  }

  public async componentDidMount() {
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;

    const [detailsResponse, tracksResponse] = await Promise.all([
      playlistService.getPlaylistDetails(ownerId, playlistId),
      playlistService.getPlaylistTracks(ownerId, playlistId)
    ]);

    if (detailsResponse.error || tracksResponse.error) {
      this.setState({ hasGetError: true });
    } else {
      this.setState({
        playlist: detailsResponse.value,
        tracks: tracksResponse.value
      });
    }
  }

  public render() {
    const playlist = this.state.playlist;
    const tracks = this.state.tracks;
    if (typeof playlist === "undefined") {
      return null;
    }

    const playlistName = playlist.name;
    const coverUrl = playlist.coverUrl;
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        {this.state.hasGetError ? (
          <p>An error occurred retrieving some part of the playlist.</p>
        ) : (
          <div>
            <Image src={coverUrl === "" ? noArtwork : coverUrl} size="medium" />
            {typeof tracks !== "undefined" && (
              <ConvertModal
                ownerId={ownerId}
                playlistId={playlistId}
                tracks={tracks}
              />
            )}
            <Tracklist tracks={tracks} />
          </div>
        )}
      </div>
    );
  }
}

export default PlaylistDetail;
