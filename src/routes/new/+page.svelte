<script lang="ts">
    import axios from 'axios'
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { page } from "$app/stores"
    import  AuthRequired  from "$components/authRequired.svelte"
 

    let githubRepo: Object|null;
    let isOwner = true
    let branches:string[] = [];
    let selectedBranch = -1;
    // idea - reject libs less than 250 LoC ??

    // console.log($page.data.bearer)


  onMount(async () => {
    console.log($page)
  //   let au = await axios.get(`https://api.github.com/user`,{
  //   headers: {
  //     'Authorization': `Bearer ${$page.data.token}`
  //   }
  // })
  //   console.warn(au);
    // let userRes = await axios.get(`https://api.github.com/users/jon-lipstate`)
    // const user = userRes.data

    // console.warn();
  });
  async function onSubmit(e:any) {
    // isOwner=true; // reset
    // const fd = new FormData(e.target);
    // const repo = fd.get('fetchRepo');

    // const response = await fetch(`https://api.github.com/repos/${repo}`);
    // const data = await response.json();
    // console.log(`GET /repos/${repo}`,data);

    // const orgRes = await fetch(`https://api.github.com/users/${$user.login}/orgs`)
    // const orgs = await orgRes.json();
    // console.log(`GET /users/${$user.login}/orgs`,orgs);

    // if (data.owner.type == "Organization") {
    //     const org = data.owner.login
    //     // /orgs/:org/public_members/:username
    //     try {
    //         const response = await fetch(`https://api.github.com/orgs/${org}/public_members/${$user.login}`);
    //         console.log(`GET /orgs/${org}/public_members/${$user.login}`, response);

    //         if (response.status == 404){
    //             isOwner=false;
    //         }
    //     } catch(e){} // nop
        
    // }

    // const _branches = await fetch(`https://api.github.com/repos/${repo}/branches`);
    // console.log(`GET /repos/${repo}/branches`, _branches);

    // branches = (await _branches.json()).map(x=>x.name)
    
    // // const tags = await fetch(`https://api.github.com/repos/${repo}/git/refs/tags`);
    // // console.log("tags",await tags.json())
    // if (isOwner) {
    //     githubRepo = data;
    //     setTimeout(()=>{updateRepo()},10); // need a tick to occur so the form shows up
    // } else{
    //     // front-end rejection - also verify on server @submit
    //     alert("you are not an owner, cannot proceed")
    // }
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

{#if $page.data.session}
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

    <h3>New Package Generation Workflow:</h3>
    <ol>
        <li><strong>User Authentication</strong>: Require github oauth</li>
        <li><strong>Display Pulldown of User/Organizations</strong>: The user selects from a dropdown of their username and organizations they are a part of.</li>
        <li><strong>Display Pulldown of Projects</strong>: Depending on the selected user/organization, display a dropdown of available projects/repositories.</li>
        <li><strong>Fetch Project Data</strong>: Once the user selects a project and clicks on the "Fetch" button, retrieve the data for the selected project.</li>
        <li><strong>Hydrate Package</strong>: Pre-fill from github response</li>
        <li><strong>Verification</strong>: Verify important contents for package creation:
            <ul>
                <li><strong>License Existence</strong>: Check that the project has a license.</li>
                <li><strong>Readme Existence</strong>: Check that the project has a readme.</li>
                <li><strong>Submitter Rights</strong>: Validate the user is the owner / authorized collaborator.</li>
                <li><strong>Tags</strong>: Ensure tags exist and are valid.</li>
                <li><strong>Pkg File Declaration</strong>: e.g. package.json is declared and properly configured.</li>
            </ul>
        </li>
        <li><strong>Form Editing</strong>: Edit Values that are mutable
            <ul>
                <li>Select Package Kind (Lib,Demo,??)</li>
                <li>License Alias? (TODO: license.key on github gives likely correct type?? eg bsd-3-clause, prefer it?)</li>
                <li>Add Topic Tags, require 1, restrict new user to use existing?</li>
            </ul>
        </li>
        <li><strong>Validation and Error Handling</strong>: Client + Server Side Verification</li>
        <li><strong>Submission</strong>: Save to db, perhaps allow for pre-published state, or save for later with a 30day expiry?</li>
        <li><strong>Confirmation</strong>: Redirect to new package details page</li>
        <li><strong>Notifications</strong>: Opt in to Notifications (? future todo)</li>
    </ol>
</main>

{:else}
  <AuthRequired/>
{/if}

<style>

</style>