import * as React from "react";
import { ILoadedPlaylist } from "../Converter";
import { Sidebar, Button } from "semantic-ui-react";

interface IConvertPanelProps {
  playlistsToConvert: ILoadedPlaylist[];
}

interface IConvertPanelState {
  visible: boolean;
}

export default class ConvertPanel extends React.Component<
  IConvertPanelProps,
  IConvertPanelState
> {
  constructor(props: IConvertPanelProps) {
    super(props);
    this.state = {
      visible: false
    };
  }

  public render() {
    const { playlistsToConvert } = this.props;
    const { visible } = this.state;
    return (
      <div>
        <Button onClick={() => this.setState({ visible: true })}>
          View Conversions
        </Button>
        <Sidebar
          animation="overlay"
          direction="bottom"
          visible={visible}
          onHide={() => this.setState({ visible: false })}
        >
          <p>AHHHHHHHHH</p>
          {playlistsToConvert.map(playlist => (
            <p key={playlist.details && playlist.details.id}>
              {playlist.details
                ? playlist.details.name
                : "Error loading playlist details"}
            </p>
          ))}
        </Sidebar>
      </div>
    );
  }
}
