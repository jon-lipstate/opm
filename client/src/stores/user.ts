import { readable, writable } from 'svelte/store';

export const isLoggedIn = writable(false);
export const isAdmin = writable(false);

export const user = writable({data:{},login:"",avatar_url:""});

//{data:{},login:"jon-lipstate",avatar_url:"https://avatars.githubusercontent.com/u/52809771?v=4"}