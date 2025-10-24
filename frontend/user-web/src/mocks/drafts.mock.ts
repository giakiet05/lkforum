export interface Draft {
  id: number;
  title: string;
  editedTime: string;
}

export interface DraftDetail {
  title: string;
  community: string;
  bodyText?: string;
  linkUrl?: string;
  tags: string[];
  tab: "text" | "images" | "link";
}

// List of drafts (for DraftsModal)
export const mockDraftsList: Draft[] = [
  { id: 1, title: "Sample1", editedTime: "edited 3 hrs ago" },
  { id: 2, title: "Sample2", editedTime: "edited 2 hrs ago" },
];

export const totalDrafts = 20;

// Draft details (for editing in CreatePostModal)
export const mockDraftsDetails: Record<number, DraftDetail> = {
  1: {
    title: "Sample1 - Amazing Discovery",
    community: "anime",
    bodyText:
      "This is the draft content for Sample1. I was watching anime and...",
    tags: ["discussion", "anime"],
    tab: "text",
  },
  2: {
    title: "Sample2 - Check this link",
    community: "3amjokes",
    linkUrl: "https://example.com/funny-joke",
    tags: ["link"],
    tab: "link",
  },
};
