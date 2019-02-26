import * as React from "react";
import { Logo } from "./Logo";
import styled from "styled-components";

interface IGreetingProps {
  loginUrl: string;
}

const Greeting: React.FunctionComponent<IGreetingProps> = ({ loginUrl }) => {
  return (
    <div>
      <header>
        <BigLogo />
      </header>
      <p>Convert your Spotify playlists to YouTube music video playlists.</p>
      <p>
        <a href={loginUrl}>Login with Spotify</a> to get started.
      </p>
    </div>
  );
};

const BigLogo = styled(Logo)`
  font-size: 6.5em;
`;

export default Greeting;
