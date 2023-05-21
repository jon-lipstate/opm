<script lang="ts">
    let githubRepo: Object|null;
    let isOwner = true
    const user = {login:"jon-lipstate"} // TODO: ADD TO STORE
    let branches:string[] = [];
    let selectedBranch = -1;

  async function onSubmit(e:any) {
    isOwner=true; // reset
    const fd = new FormData(e.target);
    const repo = fd.get('fetchRepo');
    console.log("Repo",repo,fd);
    const response = await fetch(`https://api.github.com/repos/${repo}`);
    const data = await response.json();
    // Do something with data
    console.log("github fetch",data);
    if (data.owner.type == "Organization") {
        const org = data.owner.login
        // /orgs/:org/public_members/:username
        try {
            const response = await fetch(`https://api.github.com/orgs/${org}/public_members/${user.login}`);
            if (response.status == 404){
                isOwner=false;
            }
        } catch(e){} // nop
        
    }

    const _branches = await fetch(`https://api.github.com/repos/${repo}/branches`);
    branches = (await _branches.json()).map(x=>x.name)
    
    // const tags = await fetch(`https://api.github.com/repos/${repo}/git/refs/tags`);
    // console.log("tags",await tags.json())
    if (isOwner) {
        githubRepo = data;
        setTimeout(()=>{updateRepo()},10); // need a tick to occur so the form shows up
    } else{
        // front-end rejection - also verify on server @submit
        alert("you are not an owner, cannot proceed")
    }
  }
  function updateRepo() {
    let form = document.getElementById('repoForm') as HTMLFormElement;
    //@ts-ignore
    form.repo.value = githubRepo.full_name;

  }
  function handleSelectionChange(event:Event){
    //@ts-ignore
    selectedBranch = branches.indexOf(event.target.value)
  }
</script>

<main>
    <h1>New Package</h1>
	{#if githubRepo == null}
    <form on:submit|preventDefault={onSubmit}>
        <input type="text" name="fetchRepo" id="fetchRepo" placeholder="user/project (eg odin-lang/odin)">
        <button type="submit">Fetch</button>
    </form>
    {:else}
    <button on:click={()=>githubRepo=null}>Reset</button>
    <form id="repoForm" on:submit|preventDefault={onSubmit}>
        <input type="text" name="repo" id="repo" placeholder="user/project">
        <select bind:value={branches[selectedBranch]} on:change={handleSelectionChange}>
            {#each branches as branch,index (index)}
              <option>{branch}</option>
            {/each}
          </select>
        <button type="submit">Submit</button>
    </form>
    {/if}

</main>

<style>

</style>