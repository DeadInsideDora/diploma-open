
import { UserLoginRequest, UserRegisterRequest, UserData, MapInfoUpdateRequest, ExchangeUpdateRequest } from '@/types/user';

const AUTH_API_URL = 'http://localhost:8002';
const SHOPS_API_URL = 'http://localhost:8003';

export async function loginUser(credentials: UserLoginRequest): Promise<UserData> {
  const response = await fetch(`${AUTH_API_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(credentials),
  });
  
  if (!response.ok) {
    throw new Error('Login failed');
  }
  
  return response.json();
}

export async function registerUser(userData: UserRegisterRequest): Promise<void> {
  const response = await fetch(`${AUTH_API_URL}/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(userData),
  });
  
  if (!response.ok) {
    throw new Error('Registration failed');
  }
}

export async function updateMapInfo(mapInfo: MapInfoUpdateRequest): Promise<UserData> {
  const response = await fetch(`${AUTH_API_URL}/map-info`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(mapInfo),
  });
  
  if (!response.ok) {
    throw new Error('Failed to update map info');
  }
  
  return response.json();
}

export async function updateExchange(exchangeData: ExchangeUpdateRequest): Promise<UserData> {
  const response = await fetch(`${AUTH_API_URL}/exchange`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(exchangeData),
  });
  
  if (!response.ok) {
    throw new Error('Failed to update exchange');
  }
  
  return response.json();
}

export async function updateLoyaltyCards(cards: string[]): Promise<UserData> {
  const response = await fetch(`${AUTH_API_URL}/cards`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(cards),
  });
  
  if (!response.ok) {
    throw new Error('Failed to update loyalty cards');
  }
  
  return response.json();
}

export async function getAvailableShops(): Promise<string[]> {
  const response = await fetch(`${SHOPS_API_URL}/available-shops`, {
    method: 'GET',
    credentials: 'include',
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch available shops');
  }
  
  return response.json();
}
