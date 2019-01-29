import request, { IApiResponse } from "./HttpRequest";
import { IPlaylist } from "./PlaylistService";

export interface ITrack {
  artist: string;
  title: string;
  id: string;
}

class TrackService {
  public async convert(
    ownerId: string,
    playlistId: string,
    trackId: string,
    newPlaylistId: string
  ): Promise<IApiResponse<IPlaylist>> {
    const newPlaylistQuery = newPlaylistId
      ? `?newPlaylistId=${newPlaylistId}`
      : "";

    return await request(
      `/playlists/${ownerId}/${playlistId}/tracks/${trackId}
      }${newPlaylistQuery}`,
      true,
      "POST"
    );
  }
}

const trackService = new TrackService();
export default trackService;
