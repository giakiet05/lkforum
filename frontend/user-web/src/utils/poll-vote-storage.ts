/**
 * Temporary solution to track poll votes in localStorage
 * because backend doesn't return user_vote_ids field
 */

const POLL_VOTES_KEY = 'lk_poll_votes';

interface PollVoteRecord {
  postId: string;
  optionIds: string[];
  timestamp: number;
}

export function savePollVote(postId: string, optionId: string) {
  const votes = getPollVotes();
  const existingVote = votes.find(v => v.postId === postId);
  
  if (existingVote) {
    if (!existingVote.optionIds.includes(optionId)) {
      existingVote.optionIds.push(optionId);
      existingVote.timestamp = Date.now();
    }
  } else {
    votes.push({
      postId,
      optionIds: [optionId],
      timestamp: Date.now()
    });
  }
  
  localStorage.setItem(POLL_VOTES_KEY, JSON.stringify(votes));
}

export function removePollVote(postId: string) {
  const votes = getPollVotes();
  const filtered = votes.filter(v => v.postId !== postId);
  localStorage.setItem(POLL_VOTES_KEY, JSON.stringify(filtered));
}

export function getPollVotedOptions(postId: string): string[] {
  const votes = getPollVotes();
  const vote = votes.find(v => v.postId === postId);
  return vote?.optionIds || [];
}

export function hasVotedOnPoll(postId: string): boolean {
  return getPollVotedOptions(postId).length > 0;
}

function getPollVotes(): PollVoteRecord[] {
  try {
    const stored = localStorage.getItem(POLL_VOTES_KEY);
    return stored ? JSON.parse(stored) : [];
  } catch {
    return [];
  }
}
