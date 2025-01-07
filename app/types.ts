export interface Team {
  id: number;
  name: string;
  leader: string;
}

export interface Location {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
}

export interface Stats {
  totalLocations: number;
  preachedLocations: number;
  activeTeams: number;
  teams: Team[];
}