import * as React from "react";
import { Button, Modal, Progress } from "semantic-ui-react";
import { IApiResponse } from "../../../services/HttpRequest";
import { IPlaylist } from "../../../services/PlaylistService";
import trackService, { ITrack } from "../../../services/TrackService";

interface IConvertModalProps {
  ownerId: string;
  playlistId: string;
  tracks: ITrack[];
}

interface IConvertModalState {
  convertFailures: Array<{}>;
  currentTrackTitle: string;
  loggedInYouTube: boolean;
  percentProgress: number;
  open: boolean;
  playlistUrl?: string;
}

class ConvertModal extends React.Component<
  IConvertModalProps,
  IConvertModalState
> {
  constructor(props: Readonly<IConvertModalProps>) {
    super(props);
    this.state = {
      convertFailures: [],
      currentTrackTitle: "",
      loggedInYouTube: true,
      percentProgress: 0,
      open: false,
      playlistUrl: undefined
    };
  }

  public render() {
    const loggedInYouTube = this.state.loggedInYouTube;
    const percentProgress = this.state.percentProgress;
    const convertFailures = this.state.convertFailures;
    const playlistUrl = this.state.playlistUrl;

    const content = loggedInYouTube ? (
      <div>
        <Progress
          percent={percentProgress}
          indicating={percentProgress < 100}
          success={percentProgress === 100 && convertFailures.length === 0}
          error={percentProgress === 100 && convertFailures.length > 0}
        >
          {convertFailures.length > 0 ? (
            <p>Failed to convert {this.state.convertFailures.length} tracks.</p>
          ) : null}
        </Progress>

        {playlistUrl !== "" ? (
          <p>
            Your playlist has been converted! See it{" "}
            <a href={playlistUrl}>here</a>.
          </p>
        ) : (
          <p>Converting {this.state.currentTrackTitle}...</p>
        )}
      </div>
    ) : (
      <p>
        To begin conversion, you must first{" "}
        <a
          href={
            (process.env.REACT_APP_SPOTITUBE_HOST
              ? ""
              : "http://localhost:8080") +
            "/login/youtube?returnURL=" +
            encodeURIComponent(
              window.location.pathname + window.location.search
            )
          }
        >
          login with Google/YouTube
        </a>
        .
      </p>
    );

    return (
      <Modal
        closeOnRootNodeClick={false}
        trigger={<Button primary>Convert to YouTube</Button>}
        onOpen={() => this.handleConvertClick()}
      >
        <Modal.Header>Convert Playlist to YouTube</Modal.Header>
        <Modal.Content>{content}</Modal.Content>
      </Modal>
    );
  }

  public async handleConvertClick() {
    const open = this.state.open;
    if (!open) {
      try {
        const newPlaylist = await this.convertTracks();
        this.setState({ playlistUrl: newPlaylist.url });
      } catch (err) {
        this.setState({ loggedInYouTube: false });
      }
    }

    this.setState({ open: true });
  }

  public async convertTracks(): Promise<IPlaylist> {
    const ownerId = this.props.ownerId;
    const playlistId = this.props.playlistId;
    const tracks = this.props.tracks;

    let newPlaylist = {} as IPlaylist;
    tracks.map(async (track, idx) => {
      const response = await trackService.convert(
        ownerId,
        playlistId,
        track.id,
        newPlaylist.id
      );

      this.handleConvertErrors(response, track);
      if (typeof response.value !== "undefined") {
        newPlaylist = response.value;
      }

      this.setState({
        percentProgress: ((idx + 1) / tracks.length) * 100
      });
    });

    return newPlaylist;
  }

  private handleConvertErrors<T>(
    response: IApiResponse<IPlaylist>,
    track: ITrack
  ) {
    if (response.status !== 200) {
      if (response.status === 401) {
        throw new Error("Not logged into YouTube.");
      } else {
        console.error(
          `Failed to convert ${track.title}`,
          response.error && response.error.message
        );
        this.setState({
          convertFailures: [...this.state.convertFailures, track]
        });
      }
    }

    return response;
  }
}

export default ConvertModal;
