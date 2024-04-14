export type Notification = {
  id: string;
  title: string;
  content: string;
  user_id: string;
  order_id: number;
  is_read: boolean;
  created_at: Date;
  updated_at: Date | null;
};
