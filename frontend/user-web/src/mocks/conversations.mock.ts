export interface Message {
  id: string;
  senderId: string;
  senderName: string;
  senderAvatar: string;
  content: string;
  timestamp: string;
  isRead: boolean;
  isSent: boolean; // true if sent by current user
}

export interface Conversation {
  id: string;
  userId: string;
  userName: string;
  userAvatar: string;
  lastMessage: string;
  lastMessageTime: string;
  unreadCount: number;
  isOnline: boolean;
  isMuted: boolean;
  messages: Message[];
}

export const mockConversations: Conversation[] = [
  {
    id: "1",
    userId: "cool_Pilot604",
    userName: "cool_Pilot604",
    userAvatar: "/avatar.jpg",
    lastMessage: "You: I love you too !",
    lastMessageTime: "2m",
    unreadCount: 0,
    isOnline: true,
    isMuted: false,
    messages: [
      {
        id: "m1",
        senderId: "cool_Pilot604",
        senderName: "cool_Pilot604",
        senderAvatar: "/avatar.jpg",
        content: "Hey! How are you?",
        timestamp: "10:30 AM",
        isRead: true,
        isSent: false,
      },
      {
        id: "m2",
        senderId: "me",
        senderName: "You",
        senderAvatar: "/avatar.jpg",
        content: "I'm good! How about you?",
        timestamp: "10:32 AM",
        isRead: true,
        isSent: true,
      },
      {
        id: "m3",
        senderId: "cool_Pilot604",
        senderName: "cool_Pilot604",
        senderAvatar: "/avatar.jpg",
        content: "Doing great! Want to grab coffee later?",
        timestamp: "10:35 AM",
        isRead: true,
        isSent: false,
      },
      {
        id: "m4",
        senderId: "me",
        senderName: "You",
        senderAvatar: "/avatar.jpg",
        content: "Sure! What time works for you?",
        timestamp: "10:36 AM",
        isRead: true,
        isSent: true,
      },
      {
        id: "m5",
        senderId: "cool_Pilot604",
        senderName: "cool_Pilot604",
        senderAvatar: "/avatar.jpg",
        content: "I love you !",
        timestamp: "10:38 AM",
        isRead: true,
        isSent: false,
      },
      {
        id: "m6",
        senderId: "me",
        senderName: "You",
        senderAvatar: "/avatar.jpg",
        content: "I love you too !",
        timestamp: "10:39 AM",
        isRead: true,
        isSent: true,
      },
    ],
  },
  {
    id: "2",
    userId: "cool_Annie",
    userName: "cool_Annie",
    userAvatar: "/GirlFromNowhere.jpg",
    lastMessage: "cool_Annie: i dont know man but you seem really weird lately....",
    lastMessageTime: "1h",
    unreadCount: 2,
    isOnline: false,
    isMuted: false,
    messages: [
      {
        id: "m1",
        senderId: "cool_Annie",
        senderName: "cool_Annie",
        senderAvatar: "/GirlFromNowhere.jpg",
        content: "Did you see the new update?",
        timestamp: "9:15 AM",
        isRead: true,
        isSent: false,
      },
      {
        id: "m2",
        senderId: "me",
        senderName: "You",
        senderAvatar: "/avatar.jpg",
        content: "Not yet! What's new?",
        timestamp: "9:20 AM",
        isRead: true,
        isSent: true,
      },
      {
        id: "m3",
        senderId: "cool_Annie",
        senderName: "cool_Annie",
        senderAvatar: "/GirlFromNowhere.jpg",
        content: "i dont know man but you seem really weird lately....",
        timestamp: "9:25 AM",
        isRead: false,
        isSent: false,
      },
    ],
  },
  {
    id: "3",
    userId: "cool_BigSis",
    userName: "cool_BigSis",
    userAvatar: "/avatar.jpg",
    lastMessage: "cool_BigSis: Say hi to eachother to start a conversation !",
    lastMessageTime: "3h",
    unreadCount: 0,
    isOnline: true,
    isMuted: false,
    messages: [
      {
        id: "m1",
        senderId: "cool_BigSis",
        senderName: "cool_BigSis",
        senderAvatar: "/avatar.jpg",
        content: "Say hi to eachother to start a conversation !",
        timestamp: "8:00 AM",
        isRead: true,
        isSent: false,
      },
    ],
  },
];
