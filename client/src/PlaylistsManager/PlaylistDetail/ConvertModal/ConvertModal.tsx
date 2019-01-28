import React, { Component } from "react";
import { Button, Modal, Progress } from "semantic-ui-react";
import request, { IApiResponse } from "../../../services/HttpRequest";
import trackService, { IPlaylist } from "../../../services/TrackService";

interface ITrack {
  title: string;
  id: string;
}

interface IConvertModalProps {
  ownerId: string;
  playlistId: string;
  tracks: ITrack[];
}

interface IConvertModalState {
  converted: boolean;
  convertFailures: Array<{}>;
  currentTrackTitle: string;
  loggedInYouTube: boolean;
  percentProgress: number;
  open: boolean;
  playlistUrl: string;
}

class ConvertModal extends React.Component<
  IConvertModalProps,
  IConvertModalState
> {
  constructor(props) {
    super(props);
    this.state = {
      converted: false,
      convertFailures: [],
      currentTrackTitle: "",
      loggedInYouTube: true,
      percentProgress: 0,
      open: false,
      playlistUrl: undefined
    };
  }

  public handleConvertClick() {
    const open = this.state.open;
    if (!open) {
      this.convertTracks()
        .then(res => {
          this.setState({ converted: true });
          res.json().then(playlistInfo => {
            this.setState({ playlistUrl: playlistInfo.URL });
          });
        })
        .catch(err => {
          // Catch YouTube login errors
          console.error(err);
        });
    }

    this.setState({ open: true });
  }

  public async convertTracks(): IPlaylist {
    const ownerId = this.props.ownerId;
    const playlistId = this.props.playlistId;
    const tracks = this.props.tracks;

    let newPlaylistId;

    try {
      tracks.map(async (track, idx) => {
        const response = await trackService.convert(
          ownerId,
          playlistId,
          track.id,
          newPlaylistId
        );
        this.handleConvertErrors(response, track);
        this.setState({
          percentProgress: ((idx + 1) / tracks.length) * 100
        });
      });
    } catch (err) {
      this.setState({ loggedInYouTube: false });
    }
  }

  private handleConvertErrors<T>(
    response: IApiResponse<IPlaylist>,
    track: ITrack
  ) {
    if (response.status !== 200) {
      if (response.status === 401) {
        throw new Error(response);
      } else {
        console.error(`Failed to convert ${track.title}`);
        this.setState({
          convertFailures: [...this.state.convertFailures, track]
        });
      }
    }

    return response;
  }

  public render() {
    const loggedInYouTube = this.state.loggedInYouTube;
    const percentProgress = this.state.percentProgress;
    const convertFailures = this.state.convertFailures;
    const converted = this.state.converted;
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

        {converted ? (
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
}

export default ConvertModal;
