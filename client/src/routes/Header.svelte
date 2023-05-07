<script>
	import logo from '$lib/images/odin-logo.svg';
	// import github from '$lib/images/git_light.svg';
</script>

<header>
	<nav>
		<!-- Logo-->
		<div class="logo">
			<a href="/">
				<img src={logo} alt="OdinLogo" />
			</a>
		</div>
		<!-- Search Bar-->
		<form class="search-container" action="/search">
			<input
				type="search"
				name="query"
				class="search-input"
				placeholder="Search"
				spellcheck="false"
				autocomplete="off"
			/>
			<div class="shortcut">
				<kbd class="">Ctrl</kbd> <kbd class="">K</kbd>
			</div>
		</form>
		<!-- Right Items -->
		<div class="right-items">
			<a href="https://odin-lang.org/"> ODIN-LANG </a>
			<span style="padding-left:1rem;" />
			<a href="/login"> Login </a>
		</div>
	</nav>
</header>

<style>
	header {
		display: flex;
		justify-content: space-between;
	}
	.logo img {
		margin-left: 0.75rem;
		margin-top: 0.75rem;
		width: 6rem;
		height: 2rem;
		object-fit: contain;
	}

	nav {
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 6vh;
		min-height: 55px;
		margin: 0 auto;
		background-color: var(--color-bg-0);
	}
	.search-container {
		position: relative;
		display: flex;
		align-items: center;
		width: 40rem;
		height: 100%;
		padding: 0 3px;
	}
	.search-input {
		padding: 0.2em 0.5em;
		border: 1px solid var(--back-translucent);
		font-family: inherit;
		font-size: 1.4rem;
		appearance: none;
		-webkit-appearance: none;
		width: 100%;
		height: 1.8em;
		border-radius: var(--border-radius);
	}
	.shortcut {
		/* color: var(--sk-text-3); */
		position: absolute;
		top: calc(50%-0.8rem); /*TODO: why the 0.8rem?*/
		right: 1rem;
		width: 100%;
		text-align: right;
		pointer-events: none;
		font-size: 1.2rem;
		text-transform: uppercase;
	}
	.right-items {
		margin-top: 1.2rem; /**todo: hack fix*/
		margin-right: 1rem;
		justify-content: end;
		display: flex;
		/* width: auto; */
		height: 100%;
	}

	a {
		color: var(--color-text);
	}
</style>
