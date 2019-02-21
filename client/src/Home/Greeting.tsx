import * as React from "react";
import { Header } from "semantic-ui-react";

interface IGreetingProps {
  loginUrl: string;
}

const Greeting: React.FunctionComponent<IGreetingProps> = ({ loginUrl }) => {
  return (
    <div>
      <header>
        <Header as="h1">Playlist Exchange</Header>
      </header>
      <p>Convert your Spotify playlists to YouTube music video playlists.</p>
      <p>
        <a href={loginUrl}>Login with Spotify</a> to get started.
      </p>
    </div>
  );
};

export default Greeting;
