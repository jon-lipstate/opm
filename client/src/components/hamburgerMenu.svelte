<script lang="ts">
    import { onMount } from 'svelte';
    import { isAdmin, isLoggedIn } from '$stores/user';
    
    let isOpen = false;
    function toggleMenu() {
      isOpen = !isOpen;
    }
    
    function handleLogout() {
      isLoggedIn.set(false);
    }
    
    onMount(() => {
      const closeMenu = (event:any) => {
        if (!event.target.closest('.menu') && !event.target.closest('.hamburger')) {
          isOpen = false;
        }
      };
    
      window.addEventListener('click', closeMenu);
      return () => {
        window.removeEventListener('click', closeMenu);
      };
    });
    </script>
    
    <nav>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="hamburger" on:click={toggleMenu}>
        <div class="line" class:top-line={isOpen}></div>
        <div class="line" class:mid-line={isOpen}></div>
        <div class="line" class:bot-line={isOpen}></div>
      </div>
    
      {#if isOpen}
        <div class="menu">
          <a href="/new">New Package</a>
          <a href="/manage">Manage Packages</a>
      {#if $isAdmin}
          <a href="/admin">Admin</a>
      {/if}
          <a href="/account">Account</a>
          <a href="/" on:click={handleLogout}>Logout</a>
        </div>
      {/if}
    </nav>
    
    <style>
    .hamburger {
      display: flex;
      flex-direction: column;
      justify-content: space-around;
      width: 2rem;
      height: 2rem;
      cursor: pointer;
    }
    
    .line {
      width: 2rem;
      height: 0.25rem;
      background: var(--color-theme-4);
      transition: all 0.3s ease;
    }
    
    .top-line {
        transform: translateY(0.675rem) rotate(-45deg);
    }
    
    .mid-line {
      opacity: 0;
    }
    
    .bot-line {
        transform: translateY(-0.675em) rotate(45deg);
    }
    
    .menu {
      display: flex;
      flex-direction: column;
      position: absolute;
      right: 0;
      width: 100%;
      background: var(--color-bg-1);
      text-align: center;
    }
    
    .menu a {
      padding: 1rem;
      color: var(--color-theme-4);
      text-decoration: none;
    }
    </style>
    