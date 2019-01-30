import * as React from "react";
import { Button, Modal, Progress } from "semantic-ui-react";
import { IApiResponse } from "../../../services/HttpRequest";
import { IPlaylist } from "../../../services/PlaylistService";
import trackService, { ITrack } from "../../../services/TrackService";
import { History } from "history";

interface IConvertModalProps {
  ownerId: string;
  playlistId: string;
  tracks: ITrack[];
  beginConverting: boolean;
  history: History<any>;
}

interface IConvertModalState {
  convertFailures: Array<{}>;
  currentTrackTitle: string;
  loggedInYouTube: boolean;
  percentProgress: number;
  open: boolean;
  playlistUrl?: string;
  attemptedConversion: boolean;
}

class ConvertModal extends React.Component<
  IConvertModalProps,
  IConvertModalState
> {
  constructor(props: Readonly<IConvertModalProps>) {
    super(props);
    this.state = {
      attemptedConversion: false,
      convertFailures: [],
      currentTrackTitle: "",
      loggedInYouTube: true,
      open: props.beginConverting,
      percentProgress: 0,
      playlistUrl: undefined
    };
  }

  public async componentDidMount() {
    if (this.props.beginConverting) {
      await this.attemptConvert();
    }
  }

  public render() {
    const loggedInYouTube = this.state.loggedInYouTube;
    const percentProgress = this.state.percentProgress;
    const convertFailures = this.state.convertFailures;
    const open = this.state.open;
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

        {typeof playlistUrl !== "undefined" ? (
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
      <div>
        <Button
          onClick={async () => {
            this.setState({ open: true });
            await this.attemptConvert();
          }}
          primary
        >
          Convert to YouTube
        </Button>

        <Modal
          onClose={() => {
            this.setState({ open: false });
          }}
          open={open}
        >
          <Modal.Header>Convert Playlist to YouTube</Modal.Header>
          <Modal.Content>{content}</Modal.Content>
        </Modal>
      </div>
    );
  }

  private async attemptConvert() {
    const attemptedConversion = this.state.attemptedConversion;
    if (!attemptedConversion) {
      this.props.history.push("?convert=true");

      try {
        const newPlaylist = await this.convertTracks();
        this.setState({ playlistUrl: newPlaylist.url });
      } catch (err) {
        this.setState({ loggedInYouTube: false });
      }
    }

    this.setState({ attemptedConversion: true });
  }

  private async convertTracks(): Promise<IPlaylist> {
    const ownerId = this.props.ownerId;
    const playlistId = this.props.playlistId;
    const tracks = this.props.tracks;

    let newPlaylist = {} as IPlaylist;
    for (let i = 0; i < tracks.length; i++) {
      const track = tracks[i];
      this.setState({
        currentTrackTitle: track.title
      });

      const response = await trackService.convertTrack(
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
        percentProgress: ((i + 1) / tracks.length) * 100
      });
    }

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
