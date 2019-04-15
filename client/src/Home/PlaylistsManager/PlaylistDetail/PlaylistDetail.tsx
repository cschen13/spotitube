import { parse } from "query-string";
import * as React from "react";
import { RouteComponentProps } from "react-router-dom";
import { Header, Image, Accordion, Loader } from "semantic-ui-react";
// @ts-ignore: https://github.com/Microsoft/TypeScript/issues/15146
import noArtwork from "../../../imgs/no-artwork.png";
import playlistService, { IPlaylist } from "../../../services/PlaylistService";
import { ITrack } from "../../../services/TrackService";
import ConvertModal from "./ConvertModal/ConvertModal";
import Tracklist from "./Tracklist/Tracklist";

// TODO: Pass playlist details as props instead of using API to set state
interface IPlaylistDetailState {
  readonly hasDetailsError: boolean;
  readonly hasTracksError: boolean;
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
      hasDetailsError: false,
      hasTracksError: false,
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

    const detailsResponse = await playlistService.getPlaylistDetails(
      ownerId,
      playlistId
    );
    if (detailsResponse.error) {
      this.setState({ hasDetailsError: true });
    } else {
      this.setState({
        playlist: detailsResponse.value
      });
    }

    const tracksResponse = await playlistService.getPlaylistTracks(
      ownerId,
      playlistId
    );
    if (tracksResponse.error) {
      this.setState({ hasTracksError: true });
    } else {
      this.setState({
        tracks: tracksResponse.value
      });
    }
  }

  public render() {
    const { playlist, tracks, hasDetailsError, hasTracksError } = this.state;
    if (typeof playlist === "undefined") {
      return null;
    }

    const playlistName = playlist.name;
    const coverUrl = playlist.coverUrl;
    const ownerId = this.props.match.params.ownerId;
    const playlistId = this.props.match.params.playlistId;

    const history = this.props.history;
    const beginConverting =
      parse(this.props.location.search).convert === "true";

    return (
      <div>
        <Header as="h2">{playlistName}</Header>
        {this.state.hasDetailsError && (
          <p>An error occurred retrieving some of the playlist details.</p>
        )}
        <div>
          <Image src={coverUrl === "" ? noArtwork : coverUrl} size="medium" />
          {hasTracksError && (
            <p>
              Cannot convert. An error occurred while retrieving the tracklist.
            </p>
          )}

          {typeof tracks === "undefined" ? (
            <Loader
              active
              inline="centered"
              content="Loading tracklist... (Long playlists may take a while)"
            />
          ) : (
            <div>
              <ConvertModal
                ownerId={ownerId}
                playlistId={playlistId}
                tracks={tracks}
                beginConverting={beginConverting}
                history={history}
              />
              <Accordion
                panels={[
                  {
                    content: {
                      content: <Tracklist tracks={tracks} />
                    },
                    title: {
                      content: "View tracks"
                    }
                  }
                ]}
              />
            </div>
          )}
        </div>
      </div>
    );
  }
}

export default PlaylistDetail;
