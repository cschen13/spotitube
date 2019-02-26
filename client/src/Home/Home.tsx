import * as React from "react";
import { Header } from "semantic-ui-react";
import playlistService, { IPlaylist } from "../services/PlaylistService";
import Greeting from "./Greeting";
import PlaylistsManager from "./PlaylistsManager/PlaylistsManager";
import { LogoCondensed } from "./Logo";
import styled from "styled-components";

interface IHomeState {
  loggedIn: boolean;
  playlists?: IPlaylist[];
  hasGetError: boolean;
}

class Home extends React.Component<{}, IHomeState> {
  private loginUrl =
    (process.env.REACT_APP_SPOTITUBE_HOST ? "" : "http://localhost:8080") +
    "/login/spotify?returnURL=" +
    encodeURIComponent(window.location.pathname) +
    window.location.search.replace("?", "&");

  constructor(props: Readonly<{}>) {
    super(props);
    this.state = {
      hasGetError: false,
      loggedIn: false,
      playlists: undefined
    };
  }

  public async componentDidMount() {
    const response = await playlistService.getCurrentUserPlaylists();
    const hasGetError = response.status !== 200 && response.status !== 401;
    this.setState({
      hasGetError,
      loggedIn: response.status !== 401,
      playlists: response.value
    });
  }

  public render() {
    const { hasGetError, loggedIn, playlists } = this.state;

    return loggedIn && playlists ? (
      <div>
        <header>
          <SmallLogoCondensed />
        </header>
        <PlaylistsManager playlists={playlists} hasGetError={hasGetError} />
      </div>
    ) : (
      <Greeting loginUrl={this.loginUrl} />
    );
  }
}

const SmallLogoCondensed = styled(LogoCondensed)`
  font-size: 1em;
`;

export default Home;
