import type { Philosopher } from "./types"
import { PUBLIC_GO_API_BASE } from '$env/static/public';

const API_BASE = import.meta.env.PUBLIC_GO_API_BASE || 'http://localhost:8080';
console.log('API_BASE:', API_BASE);

const PATHS = {
  philosophers: '/philosophers',
};

export const loadPhilosophers = async (): Promise<Philosopher[]> => {
  const url = `http://localhost:8080/philosophers`;

  try { 
    const res = await fetch(url);

    // Check HTTP response
    if (!res.ok) {
      console.warn(`Failed to fetch philosophers: ${res.status} ${res.statusText}`);
      return [];
    }

    const data = await res.json();
    console.log('Received data:', data);
    return data as Philosopher[];

  } catch (error) {
    console.error('Network or fetch error:', error);
    return [];
  }
};