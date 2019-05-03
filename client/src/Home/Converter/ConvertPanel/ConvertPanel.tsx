import * as React from "react";
import { ILoadedPlaylist } from "../Converter";
import { Sidebar, Button, Segment } from "semantic-ui-react";
import styled from "styled-components";

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
        <BottomBar onClick={() => this.setState({ visible: true })}>
          <p>View Conversions</p>
        </BottomBar>
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

const BottomBar = styled.div`
  position: fixed;
  left: 0;
  bottom: 0;
  height: 5vh;
  width: 100vw;
  padding: 0 10px;

  display: flex;
  flex-direction: column;
  justify-content: center;

  background-color: white;
  border-top: 1px solid;
`;
