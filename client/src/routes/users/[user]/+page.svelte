  <script lang="ts">
    import SearchResult from '../../../components/searchResult.svelte';
    import { page } from '$app/stores';
  
    export let user: string;
    export let packages: any[] = [{
			name: 'pico editor',
			version: '0.0.1',
			updated: '2/2/2902',
			downloads: 10,
			tags: ['tui', 'text-editor'],
			kind: 'demo'
		}]; 
  
    $: user = $page.params.user;
  </script>
  
  <section>
    <h1>User Profile: <span class="user">{user}</span></h1>
  
    <div>
      <h3>Packages:</h3>
      {#each packages as pkg}
      <SearchResult {pkg} />
    {/each}
    </div>
    
  </section>
  
  <style>
    .user{
      color:var(--color-theme-1);
      font-style: italic;
      font-weight: 700;
    }
  </style>