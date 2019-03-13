import * as React from "react";
import { Logo, LogoCondensed } from "./Logo";
import styled from "styled-components";
import { Header, Button } from "semantic-ui-react";

interface IGreetingProps {
  loginUrl: string;
}

const Greeting: React.FunctionComponent<IGreetingProps> = ({ loginUrl }) => {
  return (
    <div>
      <Header>
        <MobileLogo />
        <DesktopLogo />
      </Header>
      <Landing>
        <h2>
          Convert your Spotify playlists to YouTube music video playlists.
        </h2>
        <p>Get started now.</p>
        <a href={loginUrl}>Login with Spotify</a>
      </Landing>
    </div>
  );
};

const Landing = styled.div`
  height: 50vh;
  display: flex;
  flex-direction: column;
  margin-top: 50px;

  font-weight: 200;
`;

const DesktopLogo = styled(Logo)`
  margin-top: 0;
  display: none;

  font-size: 3em;

  @media only screen and (min-device-width: 768px) {
    display: block;
  }
`;

const MobileLogo = styled(LogoCondensed)`
  margin-top: 0;

  font-size: 3em;

  @media only screen and (min-device-width: 768px) {
    display: none;
  }
`;

export default Greeting;
