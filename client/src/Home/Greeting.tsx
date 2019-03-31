import * as React from "react";
import { Logo, LogoCondensed } from "./Logo";
import styled from "styled-components";
import { Header, Button } from "semantic-ui-react";

interface IGreetingProps {
  loginUrl: string;
}

const Greeting: React.FunctionComponent<IGreetingProps> = ({ loginUrl }) => {
  return (
    <Landing>
      <Header className="ui container">
        <MobileLogo />
        <DesktopLogo />
      </Header>
      <Splash>
        <SplashHeaderContainer className="ui container">
          <SplashHeader>
            Convert your Spotify playlists to YouTube music video playlists.
          </SplashHeader>
        </SplashHeaderContainer>
        <Action>
          <div className="ui container">
            <p>Link your accounts to get started now.</p>
            <LoginButton href={loginUrl} size="large" color="purple">
              Login with Spotify
            </LoginButton>
          </div>
        </Action>
      </Splash>
    </Landing>
  );
};

const Landing = styled.div`
  display: flex;
  flex-direction: column;
  height: 100vh;
`;

const Splash = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;

  font-weight: 200;
  @media only screen and (min-device-width: 768px) {
    margin-top: 25px;
  }
`;

const SplashHeaderContainer = styled.div`
  margin-bottom: 50px;
`;

const SplashHeader = styled.h2`
  font-size: 3em;

  @media only screen and (min-device-width: 768px) {
    font-size: 4em;
  }
`;

const Action = styled.div`
  flex: 1;
  padding-top: 25px;

  font-size: 2em;
  background-color: lightgray;
`;

const LoginButton = styled(Button)`
  font-size: 4em;
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
