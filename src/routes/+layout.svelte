<script>
	import { dev } from '$app/environment';
	import { inject } from '@vercel/analytics';
	inject({ mode: dev ? 'development' : 'production' });
	//
	import Header from './Header.svelte';
	import './styles.css';
	import fjord from '$lib/images/fjord.png';
</script>

<div class="app">
	<Header />
	<img id="fjord" src={fjord} alt="fjord" />
	<span class="warning">WARNING: In Development - No Data is Persisted</span>

	<main>
		<slot />
	</main>

	<footer>
		<p>
			visit <a href="https://pkg.odin-lang.org/">pkg.odin-lang.org</a> for core/vendor packages
		</p>
	</footer>
</div>

<style>
	.warning {
		color: rgba(183, 14, 14, 0.709);
		text-align: center;
		margin-top: 1rem;
		font-style: italic;
	}
	#fjord {
		opacity: 0.05;
		color: rgba(255, 0, 0, 255);
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 100%;
		text-align: center;
		pointer-events: none;
	}

	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		max-width: 64rem;
		margin: 0 auto;
		box-sizing: border-box;
	}

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 12px;
	}

	footer a {
		font-weight: bold;
	}

	@media (min-width: 480px) {
		footer {
			padding: 12px 0;
		}
	}
</style>
