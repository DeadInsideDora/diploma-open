
export interface UserLoginRequest {
  login: string;
  password: string;
}

export interface UserRegisterRequest {
  login: string;
  password: string;
  name: string;
}

export interface UserData {
  cards: string[];
  map_info: {
    point: {
      lat: number;
      lon: number;
    };
    radius: number;
  };
  name: string;
  login: string;
  exchange: number;
}

export interface MapInfoUpdateRequest {
  point: {
    lat: number;
    lon: number;
  };
  radius: number;
}

export interface ExchangeUpdateRequest {
  exchange: number;
}
