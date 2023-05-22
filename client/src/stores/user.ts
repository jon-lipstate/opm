import { readable, writable } from 'svelte/store';

export const isLoggedIn = writable(true);
export const isAdmin = writable(true);

export const user = readable({login:"jon-lipstate",avatar_url:"https://avatars.githubusercontent.com/u/52809771?v=4"});