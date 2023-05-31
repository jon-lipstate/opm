<svg width="800" height="800" viewBox="0 0 200 200">
	<path
		id="svgPath"
		d="M 8,218 150,76 82,8 14,76 156,218"
		style="
          fill: none;
          stroke: #000000;
          stroke-width: 16;
          stroke-linecap: butt;
          stroke-linejoin: miter;
          stroke-miterlimit: 4;
          stroke-dasharray: none;
        "
	/>
</svg>

<style>
	@keyframes draw {
		0% {
			stroke-dasharray: 0, 1000;
			stroke-opacity: 0;
		}
		70% {
			stroke-dasharray: 1000, 0;
			stroke-opacity: 1;
		}
		100% {
			stroke-dasharray: 1000, 0;
			stroke-opacity: 0;
		}
	}

	#svgPath {
		stroke-dasharray: 1000;
		stroke-dashoffset: 1000;
		stroke-opacity: 0;
		animation: draw 4s linear infinite;
	}
</style>
