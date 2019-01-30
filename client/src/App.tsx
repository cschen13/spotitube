import * as React from "react";
import { Header } from "semantic-ui-react";
import "./App.css";
import PlaylistsManager from "./PlaylistsManager/PlaylistsManager";
import playlistService, { IPlaylist } from "./services/PlaylistService";

interface IAppState {
  loggedIn: boolean;
  playlists?: IPlaylist[];
  hasGetError: boolean;
}

class App extends React.Component<{}, IAppState> {
  private loginUrl =
    (process.env.REACT_APP_SPOTITUBE_HOST ? "" : "http://localhost:8080") +
    "/login/spotify?returnURL=" +
    encodeURIComponent(window.location.pathname + window.location.search);

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

    let landing;
    if (hasGetError) {
      landing = (
        <p>
          An error occurred while retrieving your playlists. Try again later.
        </p>
      );
    } else if (loggedIn && playlists) {
      landing = <PlaylistsManager playlists={playlists} />;
    } else {
      landing = <Greeting loginUrl={this.loginUrl} />;
    }

    return (
      <div className="ui main container">
        <header>
          <Header as="h1">SpotiTube</Header>
        </header>
        {landing}
      </div>
    );
  }
}

interface IGreetingProps {
  loginUrl: string;
}

const Greeting: React.FunctionComponent<IGreetingProps> = ({ loginUrl }) => {
  return (
    <div>
      <p>Convert your Spotify playlists to YouTube music video playlists.</p>
      <p>
        <a href={loginUrl}>Login with Spotify</a> to get started.
      </p>
    </div>
  );
};

export default App;
