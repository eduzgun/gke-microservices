import type { Philosopher } from "./types"
import { PUBLIC_GO_API_BASE } from '$env/static/public';

const API_BASE = PUBLIC_GO_API_BASE || 'http://localhost:8080';
console.log('API_BASE:', API_BASE);

const PATHS = {
  philosophers: '/philosophers',
};

export const loadPhilosophers = async (): Promise<Philosopher[]> => {
  const url = `${API_BASE}/philosophers`;
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

export const loadPhilosopher = async (id: string | number): Promise<Philosopher | null> => {
  const url = `${API_BASE}/philosophers/${id}`;
  try {
    const res = await fetch(url);
    if (!res.ok) {
      console.warn(`Failed to fetch philosopher: ${res.status} ${res.statusText}`);
      return null;
    }
    const data = await res.json();
    console.log('Received philosopher data:', data);
    return data as Philosopher;
  } catch (error) {
    console.error('Network or fetch error:', error);
    return null;
  }
};

// export const savePhilosopher = async (philosopher: Omit<Philosopher, 'id' | 'created_at'>): Promise<Philosopher | null> => {
//   const url = `${API_BASE}/philosophers`;
//   try {
//     const res = await fetch(url, {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify(philosopher)
//     });
    
//     if (!res.ok) {
//       console.warn(`Failed to save philosopher: ${res.status} ${res.statusText}`);
//       return null;
//     }
    
//     const data = await res.json();
//     console.log('Saved philosopher:', data);
//     return data as Philosopher;
//   } catch (error) {
//     console.error('Network or save error:', error);
//     return null;
//   }
// };