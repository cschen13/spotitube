import request, { IApiResponse } from "./HttpRequest";
import { ITrack } from "./TrackService";

export interface IPlaylist {
  id: string;
  ownerId: string;
  name: string;
  url: string;
  coverUrl: string;
}

class PlaylistService {
  public async getCurrentUserPlaylists(): Promise<IApiResponse<IPlaylist[]>> {
    return await request("/playlists", true);
  }

  public async getPlaylistDetails(
    ownerId: string,
    playlistId: string
  ): Promise<IApiResponse<IPlaylist>> {
    return await request(`/playlists/${ownerId}/${playlistId}`, true);
  }

  public async getPlaylistTracks(
    ownerId: string,
    playlistId: string
  ): Promise<IApiResponse<ITrack[]>> {
    return await request(`/playlists/${ownerId}/${playlistId}/tracks`, true);
  }
}

const playlistService = new PlaylistService();
export default playlistService;
