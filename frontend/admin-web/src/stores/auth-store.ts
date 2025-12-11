import { writable } from "svelte/store";
import { tokenManager } from "../auth/token";

export const isAuthenticated = writable<boolean>(tokenManager.hasTokens());

export function setAuthenticated(value: boolean) {
  isAuthenticated.set(value);
}
