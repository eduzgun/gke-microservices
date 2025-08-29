export type Philosopher = {
  id: number;
  name: string;
  date_born: string;
  date_died: string;
  birthplace: string;
  interests: string[];
  portrait_uri: string | null;
  bio: string;
  created_at: string;  // time.Time → serialized as ISO string
};

export type Interaction = {
    id: number;
    username: string;
    type: string;
    content: string;
    created_at: string; // time.Time → serialized as ISO string
};
