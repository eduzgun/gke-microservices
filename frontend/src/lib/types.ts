export type Philosopher = {
  id: number;
  name: string;
  date_born: string;
  date_died: string;
  birthplace: string;
  interests: string[];
  portrait_uri: string;
  bio: string;
  created_at: string;  // time.Time → serialized as ISO string
};